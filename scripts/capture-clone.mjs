import { systemChromium } from "../.agents/skills/web-clone/scripts/lib/system-browser.mjs";
import fs from "node:fs";
import path from "node:path";

const targets = [
  { name: "clone-home", url: "http://localhost:5173/" },
  { name: "clone-register", url: "http://localhost:5173/register" },
  { 
    name: "clone-pdp", 
    url: "http://localhost:5173/jullindo/microsoft-surface-pro-12-copilot-pc-snapdragon-x-plus-8-core-ram-16gb-ssd-256-512gb-12-inch-touchscreen-1731633667861546606" 
  },
  { name: "clone-cart", url: "http://localhost:5173/cart" },
  { name: "clone-settings-address", url: "http://localhost:5173/user/settings/address" },
  { name: "clone-settings-payment", url: "http://localhost:5173/user/settings/payment" },
  { name: "clone-settings-bank", url: "http://localhost:5173/user/settings/bank" },
  { name: "clone-settings-mode-screen", url: "http://localhost:5173/user/settings/mode-screen" },
];

const outDir = path.resolve("docs/research/tokopedia/screenshots");
fs.mkdirSync(outDir, { recursive: true });

async function run() {
  const browser = await systemChromium.launch({ executablePath: "/usr/bin/google-chrome" });
  
  for (const target of targets) {
    console.log(`\n--- Fetching [${target.name}]: ${target.url} ---`);
    try {
      const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
      await page.goto(target.url, { timeout: 15000 });
      await new Promise((r) => setTimeout(r, 2000));

      const title = await page.evaluate(() => document.title);
      console.log(`Loaded: title="${title}"`);

      const shotPath = path.join(outDir, `${target.name}.png`);
      await page.screenshot({ path: shotPath, fullPage: false });
      console.log(`Saved screenshot: ${shotPath}`);

      await page.close().catch(() => {});
    } catch (err) {
      console.error(`Failed ${target.name}:`, err.message);
    }
  }

  await browser.close();
  console.log("\nClone capture completed!");
}

run().catch(console.error);

