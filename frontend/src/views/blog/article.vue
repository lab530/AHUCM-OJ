<template>
    <div class = "blog-card">
        <b-card class = "mt-2">
            <div class = "blog-li">
                <div class="img-box">
                    <img src="../../static/test-img/test-cover-img.png" class = "cover-img">
                </div>
                <div class = "blog-content">
                    <h3>Title</h3>
                    <p>Gratitude: a powerful practice that transforms lives. It shifts focus from lack to abundance, reduces stress, strengthens relationships, and fosters empathy. Cultivate gratitude daily for joy and a positive impact. Embrace the magic of a grateful heart.</p>
                    <hr>
                    <div class = "writer-info">
                        作者：
                        <b-avatar src="../../static/test-img/test-cover-img.png"></b-avatar>
                        &nbsp;
                        <em>zhangfei 666 </em>
                    </div>
                </div>
            </div>
        </b-card>
    </div>
</template>
<script>
export default {
  created() {
    this.idx = 0;
    // 在组件创建后立即调用的逻辑代码
    // 可以进行数据初始化、调用 API 获取数据等操作
  },
};
</script>
<style scoped>
.blog-card{
    width: 70%;
    margin: 0 auto;
    box-shadow: 1px;
}
.blog-li{
    padding: 10px;
    display: flex;
}
.cover-img{
    width: 200px;
    height: 200px;
}
.blog-content{
    width: 65;
    padding: 20px;
    position: relative;
}
.writer-info{
  position: absolute;
  right: 0;
  bottom: 0;

}
.img-box{
  width: 35%;
}

</style>