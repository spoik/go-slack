import { type RouteRecordRaw } from 'vue-router'
import Channels from '@/components/Channels.vue'

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: Channels
  },
  {
    path: '/channels/:id',
    name: 'channel',
    component: Channels
  },
]
