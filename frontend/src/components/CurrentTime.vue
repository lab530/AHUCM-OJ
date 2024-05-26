<template>
    <p>{{ date }}&nbsp;{{ hours }}:{{ minutes }}:{{ seconds }}</p>
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
.container{
    width: 180px;
}
</style>