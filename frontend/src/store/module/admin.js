import adminService from '@/service/adminService';


const adminModule = {
    namespaced: true,

    state: {
        TestCaseList: [],
        TestCaseDetail: '',
    },
    mutations: {
        setTestCaseList(state, TestCaseList) {
            state.TestCaseList = TestCaseList;
        },
        setTestCaseDetail(state, TestCaseDetail) {
            state.TestCaseDetail = TestCaseDetail;
        },
    },
    actions: {
        async fetchTestCaseList({ commit }) {
            try {
                const queryParams = window.location.search;
                // 将查询参数转换为 URL 搜索参数字符串  
                const response = await adminService.GetTestCaseList(queryParams);
                const TestCaseList = response.data;
                commit('setTestCaseList', TestCaseList);
            } catch (error) {
                // 处理请求错误  
                console.error('Request failed:', error);
                throw error;
            }
        },
        async UploadFile({ commit }, selectedFiles) {
            try {
                const queryParams = window.location.search;
                // 将查询参数转换为 URL 搜索参数字符串  
                const response = await adminService.UploadFile(queryParams, selectedFiles);
            } catch (error) {
                // 处理请求错误  
                console.error('Request failed:', error);
                throw error;
            }
        },
        async DeleteTestCase({ commit }, FileName) {
            try {
                const queryParams = window.location.search;
                // 将查询参数转换为 URL 搜索参数字符串  
                const response = await adminService.Delete(queryParams, FileName);
            } catch (error) {
                // 处理请求错误  
                console.error('Request failed:', error);
                throw error;
            }
        },
        async GetCaseInfo({ commit }) {
            try {
                const queryParams = window.location.search;
                // 将查询参数转换为 URL 搜索参数字符串  
                const response = await adminService.GetInfo(queryParams);
                return response;
            } catch (error) {
                // 处理请求错误  
                console.error('Request failed:', error);
                throw error;
            }
        },
        async UpdateCase({ commit }, content) {
            try {
                const queryParams = window.location.search;
                // 将查询参数转换为 URL 搜索参数字符串  
                const response = await adminService.Update(queryParams, content);
                return response;
            } catch (error) {
                // 处理请求错误  
                console.error('Request failed:', error);
                throw error;
            }
        },
        async GetContestDetail({ commit }) {
            try {
                const params = window.location.search;
                const response = await adminService.contestDetail(params);
                return Promise.resolve(response);
            } catch (error) {
                // 在这里处理错误，比如打印错误日志或通知用户  
                console.log(error)
                console.error('Error adding problem:', error);
                throw error; // 重新抛出错误，以便在调用处能够捕获到  
            }
        },
    },
};

export default adminModule;
