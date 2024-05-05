const problemRoutes = [
    {
        path: '/problemset',
        component: () => import('@/views/problem/problem.vue'),
        children: [
            {
                path: '',
                component: () => import('@/views/problem/problem_list.vue'),
            },
            {
                path: '/problem/status',
                component: () => import('@/views/problem/problem_status.vue'),
            },
            {
                path: '/problem/rank',
                component: () => import('@/views/problem/problem_rank.vue'),
            }
        ]
    },
    {
        path: '/problem',
        component: () => import('@/views/problem/problem_nav.vue'),
        props: (route) => ({
            pid: route.query.pid,
            cid: route.query.cid,
        }),
        children: [
            {
                path: '',
                component: () => import('@/views/problem/problem_detail.vue'),
            },
            {
                path: '/submit',
                component: () => import('@/views/problem/problem_submit.vue'),
                props: (route) => ({
                    pid: route.query.pid,
                    cid: route.query.cid,
                }),
            }
        ],

    }
];

export default problemRoutes;