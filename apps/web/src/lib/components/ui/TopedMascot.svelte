<script lang="ts">
	interface Props {
		class?: string;
		size?: 'sm' | 'md' | 'lg';
	}

	let { class: className = '', size = 'md' }: Props = $props();

	const dims = {
		sm: 'w-48 h-48',
		md: 'w-72 h-72',
		lg: 'w-96 h-96'
	};
</script>

<div class="relative flex items-center justify-center select-none {dims[size]} {className}">
	<svg viewBox="0 0 400 360" fill="none" xmlns="http://www.w3.org/2000/svg" class="h-full w-full">
		<defs>
			<radialGradient id="skyGlow" cx="50%" cy="50%" r="50%">
				<stop offset="0%" stop-color="#E8F5E9" />
				<stop offset="80%" stop-color="#E8F5E9" stop-opacity="0.5" />
				<stop offset="100%" stop-color="#FFFFFF" stop-opacity="0" />
			</radialGradient>
			<linearGradient id="storeRoof" x1="0" y1="0" x2="0" y2="1">
				<stop offset="0%" stop-color="#FF5252" />
				<stop offset="100%" stop-color="#D32F2F" />
			</linearGradient>
			<linearGradient id="storeWall" x1="0" y1="0" x2="0" y2="1">
				<stop offset="0%" stop-color="#FFF9C4" />
				<stop offset="100%" stop-color="#FFF176" />
			</linearGradient>
			<linearGradient id="boxGrad" x1="0" y1="0" x2="0" y2="1">
				<stop offset="0%" stop-color="#FFA726" />
				<stop offset="100%" stop-color="#FB8C00" />
			</linearGradient>
		</defs>

		<!-- Background Soft Glow -->
		<ellipse cx="200" cy="180" rx="180" ry="150" fill="url(#skyGlow)" />

		<!-- Parachute in Sky Left -->
		<g transform="translate(90, 80) scale(0.6)">
			<path d="M10 30 C10 5, 50 5, 50 30 Z" fill="#66BB6A" />
			<path d="M20 30 C20 12, 40 12, 40 30 Z" fill="#81C784" />
			<line x1="10" y1="30" x2="30" y2="50" stroke="#888" stroke-width="1.5" />
			<line x1="50" y1="30" x2="30" y2="50" stroke="#888" stroke-width="1.5" />
			<rect x="23" y="50" width="14" height="12" rx="2" fill="url(#boxGrad)" />
		</g>

		<!-- Parachute in Sky Right -->
		<g transform="translate(290, 50) scale(0.85)">
			<path d="M10 30 C10 5, 50 5, 50 30 Z" fill="#4CAF50" />
			<path d="M22 30 C22 10, 38 10, 38 30 Z" fill="#81C784" />
			<line x1="10" y1="30" x2="30" y2="52" stroke="#888" stroke-width="1.5" />
			<line x1="50" y1="30" x2="30" y2="52" stroke="#888" stroke-width="1.5" />
			<rect x="22" y="52" width="16" height="14" rx="2" fill="url(#boxGrad)" />
		</g>

		<!-- Store Building in Background -->
		<g transform="translate(210, 80)">
			<!-- Store Wall -->
			<rect x="0" y="70" width="130" height="140" rx="4" fill="url(#storeWall)" />
			<!-- Store Awning Canopy -->
			<path d="M-10 70 L140 70 L130 35 L0 35 Z" fill="url(#storeRoof)" />
			<!-- Canopy Stripes -->
			<path d="M15 70 L35 70 L25 35 L10 35 Z" fill="#FFFFFF" />
			<path d="M55 70 L75 70 L65 35 L50 35 Z" fill="#FFFFFF" />
			<path d="M95 70 L115 70 L105 35 L90 35 Z" fill="#FFFFFF" />
			<!-- Store Window / Door -->
			<rect x="15" y="90" width="45" height="75" rx="3" fill="#81D4FA" />
			<rect x="70" y="90" width="45" height="75" rx="3" fill="#81D4FA" />
			<line x1="37" y1="90" x2="37" y2="165" stroke="#FFFFFF" stroke-width="2" />
			<line x1="92" y1="90" x2="92" y2="165" stroke="#FFFFFF" stroke-width="2" />
		</g>

		<!-- Cardboard Box Stack Left -->
		<g transform="translate(75, 175)">
			<rect x="0" y="25" width="45" height="40" rx="3" fill="url(#boxGrad)" />
			<path d="M0 35 H45" stroke="#E65100" stroke-width="1.5" stroke-dasharray="3 2" />
			<rect x="12" y="0" width="38" height="30" rx="3" fill="#FFB74D" />
		</g>

		<!-- Main Tokopedia Owl Mascot (Front Center) -->
		<g transform="translate(130, 115)">
			<!-- Shadow underneath -->
			<ellipse cx="60" cy="180" rx="65" ry="12" fill="#B0BEC5" fill-opacity="0.4" />
			<!-- Bag Handle -->
			<path
				d="M38 50 V30 C38 18 50 10 60 10 C70 10 82 18 82 30 V50"
				stroke="#008848"
				stroke-width="9"
				stroke-linecap="round"
			/>
			<!-- Green Body Bag -->
			<rect x="10" y="45" width="100" height="115" rx="22" fill="#00AA5B" />
			<!-- Inner Bag Curved Shade -->
			<path
				d="M10 70 C25 80 50 85 60 85 C70 85 95 80 110 70 V138 C110 150 100 160 88 160 H32 C20 160 10 150 10 138 Z"
				fill="#03AC0E"
				fill-opacity="0.25"
			/>
			<!-- Left Big Eye -->
			<circle cx="42" cy="85" r="17" fill="#FFFFFF" />
			<circle cx="44" cy="85" r="10" fill="#212121" />
			<circle cx="48" cy="81" r="4" fill="#FFFFFF" />
			<!-- Right Big Eye (Winking / Smiling happy crescent) -->
			<path
				d="M68 85 C72 73 88 73 92 85"
				stroke="#FFFFFF"
				stroke-width="7"
				stroke-linecap="round"
			/>
			<path
				d="M70 85 C73 76 87 76 90 85"
				stroke="#212121"
				stroke-width="4.5"
				stroke-linecap="round"
			/>
			<!-- Yellow Beak Smile -->
			<path
				d="M60 92 L53 104 C53 107 67 107 67 104 Z"
				fill="#FFB300"
				stroke="#FF8F00"
				stroke-width="1.5"
			/>
			<!-- Rosy Cheeks -->
			<ellipse cx="28" cy="100" rx="7" ry="4" fill="#FF8A80" fill-opacity="0.6" />
			<ellipse cx="92" cy="100" rx="7" ry="4" fill="#FF8A80" fill-opacity="0.6" />
		</g>

		<!-- Cardboard Box Right Front -->
		<g transform="translate(225, 185)">
			<rect x="0" y="0" width="55" height="45" rx="4" fill="url(#boxGrad)" />
			<path d="M0 12 H55" stroke="#E65100" stroke-width="1.5" stroke-dasharray="4 2" />
			<rect x="18" y="16" width="18" height="12" rx="1.5" fill="#FFFFFF" fill-opacity="0.7" />
		</g>
	</svg>
</div>

