import request from '@/utils/request'

// 获取 lang
const GetLang = () => {
    return request.get('submit/lang');
}

// 提交代码

const Submit = (param, { UserId, Lang, ProblemId, Time, Code, File }) => {
    const formData = new FormData();
    formData.append('UserId', UserId);
    formData.append('Lang', Lang);
    formData.append('ProblemId', ProblemId);
    formData.append('Time', Time);
    formData.append('Code', Code);
    formData.append('File', File);

    return request.post('submit' + param, formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
};

const GetSubmitList = () => {
    return request.get('submit/list')
}

export default {
    GetLang,
    Submit,
    GetSubmitList,
};