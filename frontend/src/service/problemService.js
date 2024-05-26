import request from '@/utils/request'

// 获取问题

// 添加问题
const add = ({ Title, UserId, Description, Input, Output, SimpleInput, SimpleOutput, Illustrate, TimeLimit, MemoLimit }) => {
    const formData = new FormData();
    formData.append('Title', Title);
    formData.append('UserId', UserId);
    formData.append('Description', Description);
    formData.append('Input', Input);
    formData.append('Output', Output);
    formData.append('SimpleInput', SimpleInput);
    formData.append('SimpleOutput', SimpleOutput);
    formData.append('Illustrate', Illustrate);
    formData.append('TimeLimit', TimeLimit);
    formData.append('MemoLimit', MemoLimit);

    return request.post('problem/add', formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
};

// 获取问题列表
const GetProblemList = (params) => {
    // console.log(params)
    if (params !== undefined) return request.get('problem/list' + params)
    return request.get('problem/list');
};


// 获取问题详情
const GetProblem = (params) => {
    return request.get('problem' + params)
}


// 更改问题
const edit = (params, { Title, UserId, Description, Input, Output, SimpleInput, SimpleOutput, Illustrate, TimeLimit, MemoLimit }) => {
    const formData = new FormData();
    formData.append('Title', Title);
    formData.append('UserId', UserId);
    formData.append('Description', Description);
    formData.append('Input', Input);
    formData.append('Output', Output);
    formData.append('SimpleInput', SimpleInput);
    formData.append('SimpleOutput', SimpleOutput);
    formData.append('Illustrate', Illustrate);
    formData.append('TimeLimit', TimeLimit);
    formData.append('MemoLimit', MemoLimit);

    return request.put('problem/edit' + params, formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
};

const problemStatics = (params) => {
    return request.get('/problempass' + params)
}


export default {
    add,
    GetProblemList,
    GetProblem,
    edit,
    problemStatics,
};