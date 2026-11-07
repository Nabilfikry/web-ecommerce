<script lang="ts">
	import { onMount } from 'svelte';
	import './layout.css';
	import favicon from '#lib/assets/favicon.svg';
	import { auth } from '$lib/stores/auth.svelte';
	import Navbar from '$lib/components/layout/Navbar.svelte';
	import Footer from '$lib/components/layout/Footer.svelte';
	import ToastContainer from '$lib/components/ui/ToastContainer.svelte';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();

	onMount(() => {
		auth.init();
	});
</script>

<svelte:head>
	<title>SimpleShop — Modern E-Commerce Learning Platform</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="flex min-h-screen flex-col bg-slate-50 font-sans text-slate-900">
	<Navbar />

	<main class="flex flex-1 flex-col">
		{@render children()}
	</main>

	<Footer />
	<ToastContainer />
</div>
