<template>
    <div>TestCase Problem id : {{ pid }}</div>
</template>
<script>
export default {
    props: {  
        pid: {  
        type: String, // 或者 Number，取决于你的pid格式  
        required: true,  
        },
    },
    
}
</script>
<style scoped>

</style>
