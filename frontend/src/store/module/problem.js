import problemService from '@/service/problemService';


const problemModule = {
    namespaced: true,

    state: {
        problemList: [],
        problemDetail: '',
    },
    mutations: {
        setProblemList(state, problemList) {
            state.problemList = problemList;
        },
        setProblemDetail(state, problemDetail) {
            state.problemDetail = problemDetail;
        },
    },
    actions: {
        async add({ commit }, { Title, UserId, Description, Input, Output, SimpleInput, SimpleOutput, Illustrate, TimeLimit, MemoLimit }) {
            try {
                const res = await problemService.add({ Title, UserId, Description, Input, Output, SimpleInput, SimpleOutput, Illustrate, TimeLimit, MemoLimit });
                return res;
            } catch (error) {
                // 在这里处理错误，比如打印错误日志或通知用户  
                console.error('Error adding problem:', error);
                throw error; // 重新抛出错误，以便在调用处能够捕获到  
            }
        },
        async fetchProblemList({ commit }) {
            try {
                const queryParams = window.location.search;
                const response = await problemService.GetProblemList(queryParams);
                const problemList = response.data;
                commit('setProblemList', problemList);
            } catch (error) {
                console.error('Request failed:', error);
                throw error;
            }
        },
        async GetProblemDetail({ commit }) {
            try {
                const queryParams = window.location.search;
                // 将查询参数转换为 URL 搜索参数字符串  
                const response = await problemService.GetProblem(queryParams);
                const problemDetail = response.data;
                commit('setProblemDetail', problemDetail);
            } catch (error) {
                // 处理请求错误  
                console.error('Request failed:', error);
                throw error;
            }
        },
        async EditProblem({ commit }, { Title, UserId, Description, Input, Output, SimpleInput, SimpleOutput, Illustrate, TimeLimit, MemoLimit }) {
            try {
                const queryParams = window.location.search;
                const response = await problemService.edit(queryParams, { Title, UserId, Description, Input, Output, SimpleInput, SimpleOutput, Illustrate, TimeLimit, MemoLimit });
                return response;
            } catch (error) {
                // 在这里处理错误，比如打印错误日志或通知用户  
                console.error('Error adding problem:', error);
                throw error; // 重新抛出错误，以便在调用处能够捕获到  
            }
        },
    },
};

export default problemModule;
