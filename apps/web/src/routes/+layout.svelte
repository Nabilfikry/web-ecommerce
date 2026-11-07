<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import './layout.css';
	import favicon from '#lib/assets/favicon.svg';
	import { auth } from '$lib/stores/auth.svelte';
	import { cart } from '$lib/stores/cart.svelte';
	import Navbar from '$lib/components/layout/Navbar.svelte';
	import Footer from '$lib/components/layout/Footer.svelte';
	import ToastContainer from '$lib/components/ui/ToastContainer.svelte';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();

	const isAuthPage = $derived(
		page.url.pathname === '/login' || page.url.pathname === '/register'
	);

	onMount(() => {
		auth.init();
		cart.init();
	});
</script>

<svelte:head>
	<title>Situs Jual Beli Online Terlengkap, Mudah & Aman | TokoKuPedia</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="flex min-h-screen flex-col bg-[#F5F6F7] font-sans text-[#212121]">
	{#if !isAuthPage}
		<Navbar cartItemCount={cart.totalItems} />
	{/if}

	<main class="flex flex-1 flex-col">
		{@render children()}
	</main>

	{#if !isAuthPage}
		<Footer />
	{/if}
	<ToastContainer />
</div>
