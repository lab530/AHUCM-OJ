use actix_web::{post, web::Json, HttpResponse, Responder};
use log::debug;

use crate::{
    model::submission::Submission, service::executor::Executor, util::thread_pool::GLOB_THREAD_POOL,
};

#[post("/api/v1/submit")]
pub async fn submit(form: Json<Submission>) -> impl Responder {
    debug!("got submit post with form `{:?}`", form);

    let executor = Executor::new(
        form.source_path.clone(),
        form.lang.clone(),
        form.problem_id,
        form.mem_limit,
        form.time_limit,
        form.testcases_path.clone(),
        form.submission_id,
    );
    GLOB_THREAD_POOL.send_executor(executor);

    HttpResponse::Ok()
}
