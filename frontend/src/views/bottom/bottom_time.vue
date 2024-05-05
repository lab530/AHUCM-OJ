<template>
  <div class="container">
    <div class="timer">
          <h1>服务器当前时间：</h1>
      <div class="days">
        <div class="numbers">{{ date }}</div>date
      </div>
      <div class="hours">
        <div class="numbers">{{ hours }}</div>hours
      </div>
      <div class="minutes">
        <div class="numbers">{{ minutes }}</div>minutes
      </div>
      <div class="seconds">
        <div class="numbers">{{ seconds }}</div>seconds
      </div>
    </div>

  </div>
</template>

<script>
export default {
  data() {
    return {
      date:'',
      hours:'',
      minutes:'',
      seconds:'',
    };
  },
  mounted() {
    this.getServerTime();
    // 每秒更新服务器时间
    setInterval(() => {
      this.getServerTime();
    }, 1000);
  },
  methods: {
    getServerTime() {
      // 发起请求获取服务器时间
      // 这里使用示例代码，实际应根据您的服务器端设置进行调整
      const serverDate = new Date();
      const date = serverDate.getFullYear() + '-' + (serverDate.getMonth() + 1) + '-' + serverDate.getDate();
      const hours = serverDate.getHours();
      const minutes = serverDate.getMinutes();
      const seconds = serverDate.getSeconds();
      this.date = date;
      this.hours = hours;
      this.minutes = this.addZero(minutes);
      this.seconds = this.addZero(seconds);
    },
    addZero(value) {
      // 将单个数字补零
      return value < 10 ? `0${value}` : value;
    }
  }
};
</script>

<style scoped>

.container {
  margin: 0 auto;
}

.logo {
  font-size: 20px;
  font-weight: bold;
}

.timer {
  color: #183059;
  text-align: center;
  text-transform: uppercase;
  font-family: "Lato", sans-serif;
  font-size: .6em;
  letter-spacing: 5px;
}

.days,
.hours,
.minutes,
.seconds {
    display: inline-block;
    padding: 20px;
    border-radius: 5px;
}


.days {
    background: #2db8bd;
}

.hours {
    background: #F6F4F3;
    color: #183059;
}

.minutes {
    background: #276FBF;
}

.seconds {
    background: #F0A202;
}

.numbers {
    font-family: "Montserrat", sans-serif;
    color: #183059;
    font-size: 2em;
}
h1 {
    display: block;
    font-family: "Lato", sans-serif;
    text-align: center;
    font-size: 2em;
    text-transform: uppercase;
    letter-spacing: 5px;
    color: #183059;
}
</style>