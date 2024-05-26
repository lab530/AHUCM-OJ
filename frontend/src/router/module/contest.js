const contestRoutes = [
    {
        path: '/contestset',
        component: () => import('@/views/contest/contest_list.vue'),
    },
    {
        path: '/contest',
        component: () => import('@/views/contest/contest_info.vue'),
        props: (route) => ({
            cid: route.query.cid,
        }),
    },
    {
        path: '/contestdetail',
        component: () => import('@/views/contest/contest_detail.vue'),
        props: (route) => ({
            cid: route.query.cid,
        }),
    },
    {
        path: '/contestrank',
        component: () => import('@/views/contest/contest_rank.vue'),
        props: (route) => ({
            cid: route.query.cid,
        }),
    },
    {
        path: '/contestsubmit',
        component: () => import('@/views/contest/submit_record.vue'),
        props: (route) => ({
            cid: route.query.cid,
        }),
    },
    {
        path: '/conteststatics',
        component: () => import('@/views/contest/contest_static.vue'),
        props: (route) => ({
            cid: route.query.cid,
        }),
    },
];

export default contestRoutes;