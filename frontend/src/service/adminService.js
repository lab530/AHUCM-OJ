import request from '@/utils/request'
import Contest from '@/views/admin/contest.vue';

// 获取问题详情
const GetTestCaseList = (params) => {
    return request.get('admin/testcase' + params)
}

const UploadFile = (params, selectedFiles) => {
    const formData = new FormData();

    // 将选定的文件添加到 FormData
    for (let i = 0; i < selectedFiles.length; i++) {
        const file = selectedFiles[i];
        formData.append('files[]', file);
    }

    return request.post('admin/uploadcase' + params, formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
}

const Delete = (params, FileName) => {
    return request.delete('admin/deletecase' + params + "&fname=" + FileName);
}

const GetInfo = (params) => {
    return request.get('admin/casedetail' + params);
}

const Update = (params, content) => {
    const formData = new FormData();
    formData.append('content', content);
    return request.put('admin/updatacase' + params, formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
}

export default {
    GetTestCaseList,
    UploadFile,
    Delete,
    GetInfo,
    Update,
};