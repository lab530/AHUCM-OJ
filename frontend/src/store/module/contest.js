import contestService from '@/service/contestService';


const contestModule = {
    namespaced: true,

    state: {
        contestList: [],
        contestInfo: {},
        contestProblem: [],
    },
    mutations: {
        setContestList(state, contestList) {
            state.contestList = contestList;
        },
        setContestInfo(state, contestInfo) {
            state.contestInfo = contestInfo;
        },
        setContestProblem(state, contestProblem) {
            state.contestProblem = contestProblem
        },
    },
    actions: {
        async addContest({ commit }, { Title, StartAt, EndAt, Description, UserId, Public, Password, ProblemIDs, Participants }) {
            try {
                const res = await contestService.ContestAdd({ Title, StartAt, EndAt, Description, UserId, Public, Password, ProblemIDs, Participants });
                return res;
            } catch (error) {
                // 在这里处理错误，比如打印错误日志或通知用户  
                console.log(error)
                console.error('Error adding problem:', error);
                throw error; // 重新抛出错误，以便在调用处能够捕获到  
            }
        },
        async fetchContestList({ commit }) {
            try {
                const response = await contestService.GetContestList();
                commit('setContestList', response.data.data.data);
                return response.data.data.data;
            } catch (error) {
                console.error('Request failed:', error);
                throw error;
            }
        },
        async getContestInfo({ commit }) {
            try {
                const params = window.location.search;
                const response = await contestService.ContestInfo(params);
                return Promise.resolve(response);
            } catch (error) {
                console.log(error)
                console.error('Error adding problem:', error);
                throw error;
            }
        },
        async SubmitContestPassword({ commit }, { UserId, ContestPassword }) {
            try {
                const params = window.location.search;
                const response = await contestService.ContestPassword(params, { UserId, ContestPassword });
                return Promise.resolve(response);
            } catch (error) {
                // 在这里处理错误，比如打印错误日志或通知用户  
                console.log(error)
                console.error('Error adding problem:', error);
                throw error; // 重新抛出错误，以便在调用处能够捕获到  
            }
        },
        async VerityContest({ commit }, UserId) {
            try {
                const params = window.location.search;
                const response = await contestService.AuthContest(params, UserId);
                return Promise.resolve(response);
            } catch (error) {
                // 在这里处理错误，比如打印错误日志或通知用户  
                console.log(error)
                console.error('Error adding problem:', error);
                throw error; // 重新抛出错误，以便在调用处能够捕获到  
            }
        },
        async GetContestProblem({ commit }) {
            try {
                const params = window.location.search;
                const response = await contestService.contestProblem(params);
                return Promise.resolve(response);
            } catch (error) {
                // 在这里处理错误，比如打印错误日志或通知用户  
                console.log(error)
                console.error('Error adding problem:', error);
                throw error; // 重新抛出错误，以便在调用处能够捕获到  
            }
        },
        async GetContestSubmit({ commit }) {
            try {
                const params = window.location.search;
                const response = await contestService.contestSubmit(params);
                return Promise.resolve(response);
            } catch (error) {
                // 在这里处理错误，比如打印错误日志或通知用户  
                console.log(error)
                console.error('Error adding problem:', error);
                throw error; // 重新抛出错误，以便在调用处能够捕获到  
            }
        },
        async GetRankInfo({ commit }) {
            try {
                const params = window.location.search;
                const response = await contestService.contestRankInfo(params);
                return Promise.resolve(response);
            } catch (error) {
                // 在这里处理错误，比如打印错误日志或通知用户  
                console.log(error)
                console.error('Error adding problem:', error);
                throw error; // 重新抛出错误，以便在调用处能够捕获到  
            }
        },
        async editContest({ commit }, { ID, description, end_at, password, Public, start_at, title, UserId, Participants, ProblemIDs }) {
            try {
                const res = await contestService.UpdateContest({ ID, description, end_at, password, Public, start_at, title, UserId, Participants, ProblemIDs });
                return Promise.resolve(res);
            } catch (error) {
                // 在这里处理错误，比如打印错误日志或通知用户  
                console.log(error)
                console.error('Error adding problem:', error);
                throw error; // 重新抛出错误，以便在调用处能够捕获到  
            }
        },
    },
};

export default contestModule;
