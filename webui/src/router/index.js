import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', name: 'login', component: LoginView},
		{path: '/link1', name: 'home', component: HomeView},
		{path: '/link2', name: 'profile', component: ProfileView},
		
	]
})

export default router
