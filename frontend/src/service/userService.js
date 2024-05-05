import request from '@/utils/request';

// 用户注册
const register = ({ UserName, UserNickname, UserEmail, UserPassword }) => {
    const formData = new FormData();
    formData.append('UserName', UserName);
    formData.append('UserNickname', UserNickname);
    formData.append('UserEmail', UserEmail);
    formData.append('UserPassword', UserPassword);

    return request.post('auth/register', formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
};

// 用户登录
const login = ({ UserName, UserPassword }) => {
    const formData = new FormData();
    formData.append('UserName', UserName);
    formData.append('UserPassword', UserPassword);

    return request.post('auth/login', formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
};

// 用户编辑
const edit = ({ UserName, UserNickname, UserEmail, UserPassword, NewPassword, UserIcon, IconUpload }) => {
    const formData = new FormData();
    formData.append('UserName', UserName);
    formData.append('UserNickname', UserNickname);
    formData.append('UserEmail', UserEmail);
    formData.append('UserPassword', UserPassword);
    formData.append('NewPassword', NewPassword);
    formData.append('UserIcon', UserIcon);
    formData.append('IconUpload', IconUpload);
    return request.put('auth/edit', formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
}

// 获取用户信息
const info = (headers) => {
    return request.get('auth/info', {
        headers: headers,
    });
};

export default {
    register,
    info,
    login,
    edit,
};
