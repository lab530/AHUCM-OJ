import request from '@/utils/request'
// 获取问题详情
// 获取问题列表
const GetContestList = () => {
    // console.log(params)
    return request.get('/contestSet')
};

const ContestAdd = ({ Title, StartAt, EndAt, Description, UserId, Public, Password, ProblemIDs, Participants }) => {
    const formData = new FormData();
    formData.append('Title', Title);
    formData.append('UserId', UserId);
    formData.append('StartAt', StartAt);
    formData.append('EndAt', EndAt);
    formData.append('Description', Description);
    formData.append('Public', Public);
    formData.append('Password', Password);
    formData.append('ProblemList', ProblemIDs);
    formData.append('Participants', Participants);

    return request.post('contest/add', formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
};

const ContestPassword = (params, { UserId, ContestPassword }) => {
    const formData = new FormData();
    formData.append('UserId', UserId);
    formData.append('ContestPassword', ContestPassword);

    return request.post('contestverity' + params, formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
};

const AuthContest = (params, UserId) => {
    const formData = new FormData();
    formData.append('UserId', UserId);

    return request.post('authcontest' + params, formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
};

const ContestInfo = (params) => {
    return request.get('/contestInfo' + params)
}

const contestProblem = (params) => {
    return request.get('/contestproblem' + params)
}

const contestSubmit = (params) => {
    return request.get('/contestsubmit' + params)
}

const contestRankInfo = (params) => {
    return request.get('/contestrank' + params)
}


const UpdateContest = ({ ID, description, end_at, password, Public, start_at, title, UserId, Participants, ProblemIDs }) => {
    const formData = new FormData();
    formData.append('ID', ID);
    formData.append('Title', title);
    formData.append('UserId', UserId);
    formData.append('StartAt', start_at);
    formData.append('EndAt', end_at);
    formData.append('Description', description);
    formData.append('Public', Public);
    formData.append('Password', password);
    formData.append('ProblemList', ProblemIDs);
    formData.append('Participants', Participants);

    return request.put('/updateContest', formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
}

export default {
    GetContestList,
    ContestAdd,
    ContestInfo,
    ContestPassword,
    AuthContest,
    contestProblem,
    contestSubmit,
    contestRankInfo,
    UpdateContest
};