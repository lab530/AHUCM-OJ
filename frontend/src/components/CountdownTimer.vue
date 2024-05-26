<template>
    <div class="timer">
        <div class="content">
            <div class="numbers">{{ title }}{{ days }}:{{ hours }}:{{ minutes }}:{{ seconds }}</div>
        </div>
    </div>
</template>

<script>
export default {
    props: {
        title: {
            type: String,
            default: '倒计时',
        },
        targetDate: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            days: 0,
            hours: 0,
            minutes: 0,
            seconds: 0,
        };
    },
    mounted() {
        this.startCountdown();
    },
    methods: {
        startCountdown() {
        const targetDate = new Date(this.targetDate);
        const updateCountdown = () => {
            const currentDate = new Date();
            const timeDifference = targetDate - currentDate;

            if (timeDifference > 0) {
                this.days = Math.floor(timeDifference / (1000 * 60 * 60 * 24));
                this.hours = Math.floor(
                    (timeDifference % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)
                );
                this.minutes = Math.floor(
                    (timeDifference % (1000 * 60 * 60)) / (1000 * 60)
                );
                this.seconds = Math.floor((timeDifference % (1000 * 60)) / 1000);
                this.hours = String(this.hours).padStart(2, '0');
                this.minutes = String(this.minutes).padStart(2, '0');
                this.seconds = String(this.seconds).padStart(2, '0');
            } else {
                this.days = 0;
                this.hours = 0;
                this.minutes = 0;
                this.seconds = 0;
                clearInterval(countdownInterval);
            }
        };
        const countdownInterval = setInterval(updateCountdown, 1000);
        },
    },
};
</script>

<style scoped>
.timer {
  color: #183059;
  /* font-family: "Lato", sans-serif; */
  font-size: 15px;
  letter-spacing: 5px;
}


.numbers {
    /* font-family: "Montserrat", sans-serif; */
    color: #183059;
    font-size: 15px;
}

</style>