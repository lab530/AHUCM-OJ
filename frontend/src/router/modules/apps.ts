const Layout = () => import("@/layout/index.vue");

export default [
  {
    path: "/apps",
    name: "apps",
    component: () => import("@/views/apps/index.vue"),
    meta: {
      title: "应用",
      icon: "material-symbols-light:window",
      rank: 1
    },
    children: [
      {
        path: "/apps/terminal",
        name: "terminal",
        component: () => import("@/views/apps/terminal/index.vue"),
        meta: {
          title: "终端",
          icon: "material-symbols-light:terminal"
        }
      },
      {
        path: "/apps/others",
        name: "others",
        component: () => import("@/views/apps/terminal/index.vue"),
        meta: {
          title: "其他",
          icon: "material-symbols-light:other-admission-outline"
        }
      }
    ]
  },
] satisfies Array<RouteConfigsTable>;
