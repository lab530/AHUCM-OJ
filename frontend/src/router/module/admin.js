const adminRoutes = [
    {
        path: '/admin',
        component: () => import('@/views/admin/admin.vue'),
        children: [
            {
                path: '',
                component: () => import('@/views/admin/index.vue'),
            },
            {
                path: '/admin/problem',
                component: () => import('@/views/admin/problem/problem.vue'),
            },
            {
                path: '/admin/new',
                component: () => import('@/views/admin/new.vue'),
            },
            {
                path: '/admin/user',
                component: () => import('@/views/admin/user.vue'),
            },
            {
                path: '/admin/contest',
                component: () => import('@/views/admin/contest/contest_list.vue'),
            },
            {
                path: '/admin/contest/add',
                component: () => import('@/views/admin/contest/contest_add.vue'),
            },
            {
                path: '/admin/system',
                component: () => import('@/views/admin/system.vue'),
            },
            {
                path: '/admin/problem/add',
                component: () => import('@/views/problem/problem_add.vue'),
            },
            {
                path: '/admin/testcase',
                component: () => import('@/views/admin/problem/testcase/testcase_list.vue'),
                props: (route) => ({
                    pid: route.query.pid,
                }),
            },
            {
                path: '/admin/testdetail',
                component: () => import('@/views/admin/problem/testcase/testcase_detail.vue'),
            },
            {
                path: '/admin/problemdetail',
                component: () => import('@/views/admin/problem/problem_detail.vue'),
                props: (route) => ({
                    pid: route.query.pid,
                }),
            },
            {
                path: '/admin/problem/edit',
                component: () => import('@/views/problem/problem_edit.vue'),
            },
            {
                path: '/admin/contest/edit',
                component: () => import('@/views/admin/contest/contest_edit.vue'),
                props: (route) => ({
                    cid: route.query.cid,
                }),
            },
        ]
    }
];

export default adminRoutes;