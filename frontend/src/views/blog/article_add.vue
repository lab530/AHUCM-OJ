<template>
  <div class = "blog-add">
    <b-card>
      <div>
        <b-form  @submit.stop.prevent>
          <label for="feedback-user">文章标题：</label>
          <b-form-input 
            v-model="title" 
            :state="titleState" 
            id="feedback-user"
            placeholder="请输入文章标题"
            ></b-form-input>
          <b-form-invalid-feedback :state="titleState">
            文章标题不少于 4 个字符
          </b-form-invalid-feedback>
          <b-form-valid-feedback :state="titleState">
            Looks Good.
          </b-form-valid-feedback>
        </b-form>
      </div>
      <div>
        <br>
        <label for="feedback-user">文章简介：</label>
        <b-form-textarea
          id="textarea"
          v-model="introduce"
          placeholder="为了确保文章的质量，请为您的文章写一段 20 ~ 200 字的简介。"
          rows="4"
          max-rows="6"
        ></b-form-textarea>
        <b-form-invalid-feedback :state="introduceState">
            简介字数介于 20 ~ 200 个字符之间。
          </b-form-invalid-feedback>
          <b-form-valid-feedback :state="introduceState">
            Looks Good.
          </b-form-valid-feedback>
      </div>
      <div>
        <!-- Styled -->
        <br>
        <label for="feedback-user">封面图片：</label>
        <b-form-file
          v-model="file"
          :state="Boolean(file)"
          placeholder="文件大小不能超过 200 KB"
        ></b-form-file>
      </div>
      <br>
      <label for="feedback-user">文章内容</label>
      <mavon-editor v-model="value"/>
    <br>
    <div>
      <b-button block variant="primary">提交</b-button>
    </div>
    </b-card>
  </div>
</template>
<script>
 export default {
    data() {
      return {
        title: '',
        introduce:'',
        file:[],
        text:'',
        value:'',
      }
    },
    computed: {
      titleState() {
        if(this.title.length === 0) return null;
        return this.title.length > 4;
      },
      introduceState() {
        if(this.introduce.length === 0) return null;
        return this.introduce.length > 19 && this.introduce.length < 200;
      }
    }
  }
</script>

<style scoped>
.blog-add{
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

