use std::{
    sync::{
        atomic::{AtomicBool, AtomicUsize, Ordering},
        mpsc::{channel, Receiver, Sender},
        Arc, Condvar, Mutex,
    },
    thread,
};

type Thunk<'a> = Box<dyn FnOnce() + Send + 'a>;

struct SharedData {
    pub sender: Arc<Sender<Thunk<'static>>>,
    pub receiver: Mutex<Receiver<Thunk<'static>>>,

    pub empty_trigger: Mutex<()>,
    pub empty_condvar: Condvar,
    pub join_timestamp: AtomicUsize,

    pub max_thread_count: AtomicUsize,
    pub active_thread_count: AtomicUsize,
    pub panic_thread_count: AtomicUsize,
    pub queued_job_count: AtomicUsize,

    pub is_active: AtomicBool,
    pub is_done: AtomicBool,
}

impl SharedData {
    pub fn join_timestamp(&self) -> usize {
        self.join_timestamp.load(Ordering::SeqCst)
    }

    pub fn max_thread_count(&self) -> usize {
        self.max_thread_count.load(Ordering::Relaxed)
    }

    pub fn active_thread_count(&self) -> usize {
        self.active_thread_count.load(Ordering::SeqCst)
    }

    pub fn panic_thread_count(&self) -> usize {
        self.panic_thread_count.load(Ordering::SeqCst)
    }

    pub fn queued_job_count(&self) -> usize {
        self.queued_job_count.load(Ordering::SeqCst)
    }

    pub fn is_idle(&self) -> bool {
        self.queued_job_count() == 0 && self.active_thread_count() == 0
    }

    pub fn is_active(&self) -> bool {
        self.is_active.load(Ordering::SeqCst)
    }

    pub fn is_done(&self) -> bool {
        self.is_done.load(Ordering::SeqCst)
    }

    pub fn notify_when_idle(&self) {
        if self.is_idle() {
            self.empty_condvar.notify_all();
        }
    }
}

struct Sentinel<'a> {
    id: usize,
    shared_data: &'a Arc<SharedData>,
    active: bool,
}

impl<'a> Sentinel<'a> {
    pub fn new(id: usize, shared_data: &'a Arc<SharedData>) -> Self {
        Self {
            id,
            shared_data,
            active: true,
        }
    }

    pub fn cancel(&self) {
        self.active = false;
    }
}

impl<'a> Drop for Sentinel<'a> {
    fn drop(&mut self) {
        if self.active {
            self.shared_data
                .active_thread_count
                .fetch_sub(1, Ordering::SeqCst);
            if thread::panicking() {
                self.shared_data
                    .panic_thread_count
                    .fetch_add(1, Ordering::SeqCst);
            }
            ThreadPool::spawn_thread(self.id, self.shared_data.clone());
            self.shared_data.notify_when_idle();
        }
    }
}

pub struct ThreadPool {
    pub sender: Arc<Sender<Thunk<'static>>>,
    pub shared_data: Arc<SharedData>,
}

impl ThreadPool {
    pub fn new(size: usize) -> Self {
        let (sender, receiver) = channel::<Thunk<'static>>();
        let sender = Arc::new(sender);
        let shared_data = Arc::new(SharedData {
            sender: sender.clone(),
            receiver: Mutex::new(receiver),
            empty_trigger: Mutex::new(()),
            empty_condvar: Condvar::new(),
            join_timestamp: AtomicUsize::new(0),
            max_thread_count: AtomicUsize::new(size),
            active_thread_count: AtomicUsize::new(0),
            panic_thread_count: AtomicUsize::new(0),
            queued_job_count: AtomicUsize::new(0),
            is_active: AtomicBool::new(false),
            is_done: AtomicBool::new(false),
        });

        Self {
            sender,
            shared_data,
        }
    }

    pub fn max_thread_count(&self) -> usize {
        self.shared_data.max_thread_count()
    }

    pub fn active_thread_count(&self) -> usize {
        self.shared_data.active_thread_count()
    }

    pub fn panic_thread_count(&self) -> usize {
        self.shared_data.panic_thread_count()
    }

    pub fn queued_job_count(&self) -> usize {
        self.shared_data.queued_job_count()
    }

    pub fn is_idle(&self) -> bool {
        self.shared_data.is_idle()
    }

    pub fn awake_all(&self) {
        self.shared_data.is_active.store(true, Ordering::SeqCst);
    }

    pub fn block_all(&self) {
        self.shared_data.is_active.store(false, Ordering::SeqCst);
    }

    pub fn stop_all(&self, unconditional: bool) {
        if unconditional {
            let lock = self.shared_data.receiver.lock().unwrap();
            while lock.recv().is_ok() {}
        }
        self.shared_data.is_done.store(true, Ordering::SeqCst);
    }

    pub fn join(&self) {
        if self.is_idle() {
            return;
        }

        let join_timestamp = self.shared_data.join_timestamp();
        let mut trigger = self.shared_data.empty_trigger.lock().unwrap();

        while join_timestamp == self.shared_data.join_timestamp() && !self.is_idle() {
            trigger = self.shared_data.empty_condvar.wait(trigger).unwrap();
        }

        let _ = self.shared_data.join_timestamp.compare_exchange(
            join_timestamp,
            join_timestamp + 1,
            Ordering::SeqCst,
            Ordering::SeqCst,
        );
    }

    pub fn send_job<F>(&self, job: F)
    where
        F: FnOnce() + Send + 'static,
    {
        self.shared_data
            .queued_job_count
            .fetch_add(1, Ordering::SeqCst);
        self.sender.send(Box::new(job)).unwrap();
    }

    pub fn spawn_thread(id: usize, shared_data: Arc<SharedData>) {
        let builder = thread::Builder::new();

        builder
            .spawn(move || {
                let sentinel = Sentinel::new(id, &shared_data);

                loop {
                    if shared_data.is_done() {
                        break;
                    }
                    if !shared_data.is_active() {
                        break;
                    }

                    let message = {
                        let lock = shared_data.receiver.lock().unwrap();
                        lock.recv()
                    };

                    let job = match message {
                        Ok(job) => job,
                        Err(_) => continue,
                    };

                    shared_data
                        .active_thread_count
                        .fetch_add(1, Ordering::SeqCst);
                    shared_data.queued_job_count.fetch_sub(1, Ordering::SeqCst);

                    job();

                    shared_data
                        .active_thread_count
                        .fetch_sub(1, Ordering::SeqCst);
                    shared_data.notify_when_idle();
                }

                sentinel.cancel();
            })
            .unwrap();
    }
}
