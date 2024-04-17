import storageService from '@/service/storageService';

import userService from '@/service/userService';

const userModule = {
    namespaced: true,
    state: {
        token: storageService.get(storageService.USER_TOKEN),
        userInfo: null,
    },

    mutations: {
        SET_TOKEN(state, token) {
            // 更新本地缓存
            storageService.set(storageService.USER_TOKEN, token);
            // 更新 state
            state.token = token;
        },
        SET_USERINFO(state, userInfo) {
            // 更新 state
            state.userInfo = userInfo;
        },
    },

    actions: {
        register(context, { UserName, UserNickname, UserEmail, UserPassword }) {
            return new Promise((resolve, reject) => {
                userService.register({ UserName, UserNickname, UserEmail, UserPassword }).then((res) => {
                    // 保存 token
                    context.commit('SET_TOKEN', res.data.data.token);
                    return userService.info();
                }).then((res) => {
                    // 保存用户信息
                    context.commit('SET_USERINFO', res.data.data.user);
                    resolve(res);
                }).catch((err) => {
                    reject(err);
                });
            });
        },
        login(context, { UserName, UserPassword }) {
            return new Promise((resolve, reject) => {
                userService.login({ UserName, UserPassword }).then((res) => {
                    // 保存 token
                    context.commit('SET_TOKEN', res.data.data.token);
                    return userService.info();
                }).then((res) => {
                    // 保存用户信息
                    context.commit('SET_USERINFO', res.data.data.user);
                    resolve(res);
                }).catch((err) => {
                    reject(err);
                });
            });
        },
        logout({ commit }) {
            // 清除 token
            commit('SET_TOKEN', '');
            storageService.set(storageService.USER_TOKEN, '')

            window.location.reload();
        },
        edit(context, { UserName, UserNickname, UserEmail, UserPassword, NewPassword, UserIcon, IconUpload }) {
            return new Promise((resolve, reject) => {
                console.log('userEdit method called');
                console.log(Response)
                if (UserIcon) IconUpload = "YES";
                userService.edit({ UserName, UserNickname, UserEmail, UserPassword, NewPassword, UserIcon, IconUpload }).then((res) => {
                    // 清除 token
                    // context.commit('SET_TOKEN', '');
                    // storageService.set(storageService.USER_TOKEN, '');
                    // 清除用户信息
                    context.commit('SET_USERINFO', '');
                    storageService.set(storageService.USER_INFO, '');
                    // // 保存 token
                    return userService.info();
                    // 保存用户信息
                }).then((res) => {
                    context.commit('SET_USERINFO', res.data.data.user);
                    resolve(res);
                }).catch((err) => {
                    commit('SET_TOKEN', '');
                    storageService.set(storageService.USER_TOKEN, '')

                    window.location.reload();
                    reject(err);
                });
            });
        },
        async getInfo(context, state) {
            const { token } = state;
            try {
                const headers = {
                    // 在请求头中添加 token
                    Authorization: `Bearer ${token}`,
                };
                const response = await userService.info(headers);
                context.commit('SET_USERINFO', response.data.data.user);
            } catch (error) {
                // 处理请求错误  
                console.error('Request failed:', error);
                throw error;
            }
        }
    },
};

export default userModule;
