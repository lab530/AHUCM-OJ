import submitService from '@/service/submitService';

const submitModule = {
    namespaced: true,

    state: {
        languages: [],
        SubmitList: [],
    },
    mutations: {
        setLang(state, lang) {
            state.languages = lang;
        },
        setSubmitList(state, SubmitList) {
            state.SubmitList = SubmitList
        },
    },
    actions: {
        async GetLang({ commit }) {
            try {
                const response = await submitService.GetLang();
                // console.log(response)
                const lang = response.data.data.lang;
                commit('setLang', lang);
                return response;
            } catch (error) {
                // 在这里处理错误，比如打印错误日志或通知用户  
                console.error('Error get lang', error);
                throw error; // 重新抛出错误，以便在调用处能够捕获到  
            }
        },
        async fetchSubmit({ commit }, { UserId, Lang, ProblemId, Time, Code, File }) {
            try {
                const queryParams = window.location.search;
                const response = await submitService.Submit(queryParams, { UserId, Lang, ProblemId, Time, Code, File });
                return response;
            } catch (error) {
                // 在这里处理错误，比如打印错误日志或通知用户  
                console.error('Error adding problem:', error);
                throw error; // 重新抛出错误，以便在调用处能够捕获到  
            }
        },
        async fetchSubmitList({ commit }) {
            try {
                const response = await submitService.GetSubmitList();
                commit('setSubmitList', response.data.data)
            } catch (error) {
                // 在这里处理错误，比如打印错误日志或通知用户  
                console.error('Error adding problem:', error);
                throw error; // 重新抛出错误，以便在调用处能够捕获到  
            }
        }
    },

};

export default submitModule;
