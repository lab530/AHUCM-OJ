<template>
  <div class = "article">
    <div class = "problem_detail">
      <div v-if="problem && problem.data">
          <div class = "title">
            <h3>{{ problem.data.title }}</h3>
          </div>
          <div class = "limit">
              Time:&nbsp;{{ problem.data.time_limit }}&nbsp;ms &nbsp; &nbsp; Memo:&nbsp;{{ problem.data.memo_limit }}&nbsp;ms
          </div>
          <div class = "Info">最近更新: {{ localTime }} &nbsp; &nbsp; &nbsp; 命题人: &nbsp;<a href="#">{{ problem.Info.user_name }}</a></div>
      </div>
      <div v-else>
        loading...
      </div>
      <div v-if="problem && problem.data">
        <h3 class = "Header">题目描述：</h3>
          <mavon-editor 
            v-model="problem.data.description"
            class="content-show"   
            :subfield="false" 
            defaultOpen="preview"
            :editable="false" 
            :toolbarsFlag="false"
            :boxShadow="false" previewBackground="#ffffff "/>
          <h3 class = "Header">输入：</h3>
          <mavon-editor 
            v-model="problem.data.input"
            class="content-show"   
            :subfield="false" 
            defaultOpen="preview"
            :editable="false" 
            :toolbarsFlag="false"
            :boxShadow="false" previewBackground="#ffffff "/>
          <h3 class = "Header">输出：</h3>
          <mavon-editor 
            v-model="problem.data.output"
            class="content-show"   
            :subfield="false" 
            defaultOpen="preview"
            :editable="false" 
            :toolbarsFlag="false"
            :boxShadow="false" previewBackground="#ffffff "/>
          <div class = "inline">
            <h3 class = "Header">样例输入：</h3> 
            <b-button @click="copyData('input')" variant="outline-info" class="btn-sm"><b-icon icon="files"></b-icon> 复制</b-button>
          </div> 
          <b-alert show class="bg" >
            <mavon-editor 
            v-model="problem.data.simple_input"
            class="content-show"   
            :subfield="false" 
            defaultOpen="preview"
            :editable="false" 
            :toolbarsFlag="false"
            :boxShadow="false" previewBackground="#ffffff "/>
          </b-alert>
          <div class = "inline">
            <h3 class = "Header">样例输出：</h3> 
            <b-button @click="copyData('output')" variant="outline-info" class="btn-sm"><b-icon icon="files"></b-icon> 复制</b-button>
          </div>
            <b-alert show class = "bg">
            <mavon-editor 
            v-model="problem.data.simple_output"
            class="content-show"   
            :subfield="false" 
            defaultOpen="preview"
            :editable="false" 
            :toolbarsFlag="false"
            :boxShadow="false" previewBackground="#ffffff "/>
          </b-alert>
          <h3  class = "Header">数据范围与提示：</h3>
          <mavon-editor
            v-model="problem.data.illustrate"
            class="content-show"   
            :subfield="false" 
            defaultOpen="preview"
            :editable="false" 
            :toolbarsFlag="false"
            :boxShadow="false" previewBackground="#ffffff "/>
      </div>
      <div v-else>
        loading...
      </div>
    </div>
    <div class = "toolCard">
        <div ref="chart" style="width: 280px; height: 280px;"></div>
    </div>
  </div>
</template>
<script>
import { mapState, mapActions } from 'vuex';
import * as echarts from 'echarts'
export default {
data() {
  return {
      problem: '',
      localTime:'',
      Testdata: [
            { value: 1048, name: 'Search Engine' },
            { value: 735, name: 'Direct' },
            { value: 580, name: 'Email' },
            { value: 484, name: 'Union Ads' },
            { value: 300, name: 'Video Ads' }
          ],
  };
},
computed: {  
  // 使用计算属性获取 Vuex 状态中的 problemDetail  
  ...mapState('problemModule', ['problemDetail']),  
},
created() {
  this.Detail().then(() => {
    // 在获取问题列表数据后，更新 data 数组
    this.problem = this.problemDetail.data;
    if(this.problem) this.formatTimeString(this.problem.data.UpdatedAt)
  }).catch((error) => {
    this.$bvToast.toast(error.response.data.msg, {
      title: '数据验证错误',
      variant: 'danger',
      toaster:'b-toaster-bottom-right',
      solid: true,
      appendToast:true,
    });
  });
  this.SubmitStatics().then((response) => {
    // 在获取问题列表数据后，更新 data 数组
    this.initChart(response.data.data.data);
  }).catch((error) => {
    this.$bvToast.toast(error.response.data.msg, {
      title: '数据验证错误',
      variant: 'danger',
      toaster:'b-toaster-bottom-right',
      solid: true,
      appendToast:true,
    });
  });
},
methods: {
  ...mapActions('problemModule', { Detail:'GetProblemDetail'}),
  ...mapActions('problemModule',['SubmitStatics']),
  formatTimeString(GetTime) {
    const date = new Date(GetTime);
    this.localTime = date.toLocaleString(); // 使用 toLocaleString 方法根据用户的本地时间格式进行格式化
  },
  copyData(form) {
    console.log(form)
    let content = this.problem.data.simple_input;
    if(form !== 'input') content = this.problem.data.simple_output;
    navigator.clipboard.writeText(content)
      .then(() => {
        this.$bvToast.toast('内容已复制', {
          title: '消息通知',
          variant: 'success',
          toaster:'b-toaster-bottom-right',
          solid: true,
          appendToast:true,
        });
      })
      .catch(() => {
        this.$bvToast.toast('复制错误', {
          title: '消息通知',
          variant: 'danger',
          toaster:'b-toaster-bottom-right',
          solid: true,
          appendToast:true,
      });
    });
  },
  initChart(inputData) {
    const result = [];  
    for (const key in inputData) {  
      if (Object.prototype.hasOwnProperty.call(inputData, key)) {  
        const item = inputData[key]; // 获取当前对象的值  
        result.push({  
          name: item.value, // 从对象中提取 name 属性  
          value: item.name // 从对象中提取 value 属性  
        });  
      }  
    }       
    // 基于准备好的dom，初始化echarts实例
    const myChart = echarts.init(this.$refs.chart)
    // 绘制饼图
    myChart.setOption({
      title: {
        text: '题目提交统计图',
        left: 'center'
      },
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b} : {c} ({d}%)'
      },
      series: [
        {
          name: 'Problem',
          type: 'pie',
          radius: '50%',
          data: result,
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          }
        }
      ]
    })
  },
},
}
</script>

<style scoped>
.article{
  margin-left: 22px;
  display: flex;
}

.problem_detail{
width: 850px;
}
.toolCard{
margin-left: 20px;
width: 200px;
}
.title, .limit, .Info{
display: flex;
justify-content: center;
text-decoration: none;
padding: 2px;
}
.Header{
margin-top: 3%;
font-size: 18px; /* 根据需要设置合适的字体大小 */
font-weight: bold;
}
.content-show {
font-size: 16px;
min-height: 1px;
line-height: 1.8;
/* height: 100% !important; */
border: none !important;
}
/deep/ .v-note-wrapper .v-note-panel .v-note-show .v-show-content,
/deep/ .v-note-wrapper .v-note-panel .v-note-show .v-show-content-html {
height: 100%;
padding: 0;
overflow-y: auto;
box-sizing: border-box;
overflow-x: hidden;
}

/deep/ .markdown-body p{
  margin-top: 0;
  margin-bottom: 0;
}
.bg {
background-color: #ffffff;
}
.inline{
display: flex;
position: relative;
/* align-items: flex-end; */
}
.btn-sm{
position: absolute;
bottom: 0;
right: 0;
/* margin-left: 90px; */
margin-bottom: 5px;
}
</style>