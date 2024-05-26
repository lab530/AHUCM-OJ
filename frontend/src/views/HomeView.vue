<template>
  <b-card class="box text-center">
    <b-card class="img">
      <img src="../../../backend/static/carousel/4.png">
    </b-card>
  </b-card>
</template>

<script>
  export default {
    data() {
      return {
        slide: 0,
        sliding: null
      }
    },
    methods: {
      onSlideStart(slide) {
        this.sliding = true
      },
      onSlideEnd(slide) {
        this.sliding = false
      }
    }
  }
</script>

<style scoped>
.box {
  width: 1200px;
  margin: 10px auto;
  background-image: url("../../public/bg.png");
  background-repeat: repeat;
  background-size: 150px 150px; /* 自定义宽度为800像素，高度为400像素 */
}
.img{
  width: 1000px;
  margin: 0 auto;
}
.img img{
  width: 700px;
}
</style>