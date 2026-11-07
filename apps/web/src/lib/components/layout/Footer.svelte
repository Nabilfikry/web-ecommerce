<script lang="ts">
	import TokoKuPediaLogo from '../ui/TokoKuPediaLogo.svelte';
	import { ShieldCheck, Truck, Headphones, Smartphone } from 'lucide-svelte';

	const currentYear = new Date().getFullYear();
</script>

<footer class="mt-auto border-t border-[#E5E7E9] bg-white text-[#212121]">
	<!-- Value Proposition Badges Bar -->
	<div class="border-b border-[#E5E7E9] bg-[#F5F6F7] py-6">
		<div class="mx-auto grid max-w-7xl grid-cols-2 gap-4 px-4 sm:grid-cols-4 sm:px-6 lg:px-8">
			<div class="flex items-center gap-3">
				<div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-[#E8F5E9] text-[#00AA5B]">
					<ShieldCheck class="h-5 w-5" />
				</div>
				<div>
					<div class="text-xs font-bold text-[#212121]">100% Transaksi Aman</div>
					<div class="text-[11px] text-[#6D7588]">Garansi uang kembali</div>
				</div>
			</div>

			<div class="flex items-center gap-3">
				<div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-[#E8F5E9] text-[#00AA5B]">
					<Truck class="h-5 w-5" />
				</div>
				<div>
					<div class="text-xs font-bold text-[#212121]">Bebas Ongkir</div>
					<div class="text-[11px] text-[#6D7588]">Ke seluruh Indonesia</div>
				</div>
			</div>

			<div class="flex items-center gap-3">
				<div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-[#E8F5E9] text-[#00AA5B]">
					<Headphones class="h-5 w-5" />
				</div>
				<div>
					<div class="text-xs font-bold text-[#212121]">Layanan 24/7</div>
					<div class="text-[11px] text-[#6D7588]">Bantuan TokoKuPedia Care</div>
				</div>
			</div>

			<div class="flex items-center gap-3">
				<div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-[#E8F5E9] text-[#00AA5B]">
					<Smartphone class="h-5 w-5" />
				</div>
				<div>
					<div class="text-xs font-bold text-[#212121]">Belanja di Aplikasi</div>
					<div class="text-[11px] text-[#6D7588]">Lebih banyak promo spesial</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Main Footer Columns -->
	<div class="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
		<div class="grid grid-cols-2 gap-8 md:grid-cols-5">
			<!-- Col 1 -->
			<div>
				<h4 class="mb-3 text-xs font-bold tracking-wider text-[#212121] uppercase">TokoKuPedia</h4>
				<ul class="space-y-2 text-xs text-[#6D7588]">
					<li><a href="/#about" class="hover:text-[#00AA5B]">Tentang TokoKuPedia</a></li>
					<li><a href="/#karir" class="hover:text-[#00AA5B]">Karir</a></li>
					<li><a href="/#blog" class="hover:text-[#00AA5B]">Blog TokoKuPedia</a></li>
					<li><a href="/#mitra" class="hover:text-[#00AA5B]">Mitra Blog</a></li>
					<li><a href="/#official" class="hover:text-[#00AA5B]">Official Store</a></li>
				</ul>
			</div>

			<!-- Col 2 -->
			<div>
				<h4 class="mb-3 text-xs font-bold tracking-wider text-[#212121] uppercase">Beli</h4>
				<ul class="space-y-2 text-xs text-[#6D7588]">
					<li><a href="/#topup" class="hover:text-[#00AA5B]">Tagihan & Top Up</a></li>
					<li><a href="/products" class="hover:text-[#00AA5B]">Katalog Belanja</a></li>
					<li><a href="/#bebas-ongkir" class="hover:text-[#00AA5B]">Bebas Ongkir</a></li>
					<li><a href="/#cod" class="hover:text-[#00AA5B]">TokoKuPedia COD</a></li>
				</ul>
			</div>

			<!-- Col 3 -->
			<div>
				<h4 class="mb-3 text-xs font-bold tracking-wider text-[#212121] uppercase">Jual</h4>
				<ul class="space-y-2 text-xs text-[#6D7588]">
					<li><a href="/register?role=seller" class="hover:text-[#00AA5B]">Pusat Edukasi Seller</a></li>
					<li><a href="/dashboard" class="hover:text-[#00AA5B]">Dashboard Mitra Toko</a></li>
					<li><a href="/register?role=seller" class="hover:text-[#00AA5B]">Daftar Official Store</a></li>
				</ul>
			</div>

			<!-- Col 4 -->
			<div>
				<h4 class="mb-3 text-xs font-bold tracking-wider text-[#212121] uppercase">Bantuan & Panduan</h4>
				<ul class="space-y-2 text-xs text-[#6D7588]">
					<li><a href="/#care" class="hover:text-[#00AA5B]">TokoKuPedia Care</a></li>
					<li><a href="/#syarat" class="hover:text-[#00AA5B]">Syarat & Ketentuan</a></li>
					<li><a href="/#privasi" class="hover:text-[#00AA5B]">Kebijakan Privasi</a></li>
					<li><a href="/user/settings/mode-screen" class="hover:text-[#00AA5B]">Pengaturan Tampilan</a></li>
				</ul>
			</div>

			<!-- Col 5: Security & App -->
			<div class="col-span-2 md:col-span-1">
				<h4 class="mb-3 text-xs font-bold tracking-wider text-[#212121] uppercase">Keamanan & Mitra</h4>
				<div class="flex flex-wrap gap-2 text-xs text-[#6D7588]">
					<span class="rounded border border-[#E5E7E9] bg-slate-50 px-2 py-1 font-semibold text-[11px]">BCA</span>
					<span class="rounded border border-[#E5E7E9] bg-slate-50 px-2 py-1 font-semibold text-[11px]">Mandiri</span>
					<span class="rounded border border-[#E5E7E9] bg-slate-50 px-2 py-1 font-semibold text-[11px]">BRI</span>
					<span class="rounded border border-[#E5E7E9] bg-slate-50 px-2 py-1 font-semibold text-[11px]">BNI</span>
					<span class="rounded border border-[#E5E7E9] bg-slate-50 px-2 py-1 font-semibold text-[11px]">GoPay</span>
					<span class="rounded border border-[#E5E7E9] bg-slate-50 px-2 py-1 font-semibold text-[11px]">Visa</span>
					<span class="rounded border border-[#E5E7E9] bg-slate-50 px-2 py-1 font-semibold text-[11px]">Mastercard</span>
					<span class="rounded border border-[#E5E7E9] bg-slate-50 px-2 py-1 font-semibold text-[11px]">QRIS</span>
				</div>

				<div class="mt-5">
					<h4 class="mb-2 text-xs font-bold text-[#212121]">Ikuti Kami</h4>
					<div class="flex items-center gap-2 text-xs text-[#6D7588]">
						<span class="cursor-pointer hover:text-[#00AA5B]">Facebook</span> •
						<span class="cursor-pointer hover:text-[#00AA5B]">Instagram</span> •
						<span class="cursor-pointer hover:text-[#00AA5B]">Twitter</span>
					</div>
				</div>
			</div>
		</div>

		<!-- Bottom Copyright Bar -->
		<div class="mt-12 flex flex-col items-center justify-between gap-4 border-t border-[#E5E7E9] pt-6 sm:flex-row text-xs text-[#6D7588]">
			<div class="flex items-center gap-2">
				<TokoKuPediaLogo size="sm" />
				<span>&copy; 2009 - {currentYear}, PT TokoKuPedia Indonesia.</span>
			</div>
			<div class="flex items-center gap-4 text-[11px]">
				<span>Kenyamanan dan Keamanan Belanja Prioritas Kami</span>
			</div>
		</div>
	</div>
</footer>
