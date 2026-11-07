import { systemChromium } from "../.agents/skills/web-clone/scripts/lib/system-browser.mjs";
import fs from "node:fs";
import path from "node:path";

const targets = [
  { name: "home", url: "https://www.tokopedia.com/" },
  { name: "register", url: "https://www.tokopedia.com/register" },
  { 
    name: "pdp", 
    url: "https://www.tokopedia.com/jullindo/microsoft-surface-pro-12-copilot-pc-snapdragon-x-plus-8-core-ram-16gb-ssd-256-512gb-12-inch-touchscreen-1731633667861546606?t_id=1794033000856&t_st=1&t_pp=homepage&t_efo=pure_goods_card&t_ef=homepage&t_sm=rec_homepage_outer_flow&t_spt=homepage" 
  },
  { 
    name: "cart", 
    url: "https://www.tokopedia.com/cart?t_id=1794033050628&t_st=1&t_pp=homepage&t_efo=pure_goods_card&t_ef=&t_sm=rec_homepage_outer_flow&t_spt=homepage" 
  },
  { name: "settings-address", url: "https://www.tokopedia.com/user/settings/address" },
  { name: "settings-payment", url: "https://www.tokopedia.com/user/settings/payment" },
  { name: "settings-bank", url: "https://www.tokopedia.com/user/settings/bank" },
  { name: "settings-mode-screen", url: "https://www.tokopedia.com/user/settings/mode-screen" },
];

const outDir = path.resolve("docs/research/tokopedia/screenshots");
fs.mkdirSync(outDir, { recursive: true });

async function run() {
  const browser = await systemChromium.launch({ executablePath: "/usr/bin/google-chrome" });
  
  for (const target of targets) {
    console.log(`\n--- Fetching [${target.name}]: ${target.url} ---`);
    try {
      const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
      await page.connection.send("Network.setUserAgentOverride", {
        userAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
      });

      await page.goto(target.url, { timeout: 35000 });
      // wait a bit for rendering
      await new Promise((r) => setTimeout(r, 4000));

      const title = await page.evaluate(() => document.title);
      const currentUrl = await page.evaluate(() => window.location.href);
      console.log(`Loaded: title="${title}", currentUrl="${currentUrl}"`);

      // Screenshot viewport
      const shotPath = path.join(outDir, `${target.name}.png`);
      await page.screenshot({ path: shotPath, fullPage: false });
      console.log(`Saved screenshot: ${shotPath}`);

      // Extract some design tokens and dom structure
      const info = await page.evaluate(() => {
        const bodyStyles = window.getComputedStyle(document.body);
        return {
          title: document.title,
          url: window.location.href,
          backgroundColor: bodyStyles.backgroundColor,
          color: bodyStyles.color,
          fontFamily: bodyStyles.fontFamily,
          h1: Array.from(document.querySelectorAll("h1")).map((el) => el.innerText.trim()).filter(Boolean),
          buttons: Array.from(document.querySelectorAll("button")).map((el) => el.innerText.trim()).filter(Boolean).slice(0, 15),
          headerExists: !!document.querySelector("header"),
          footerExists: !!document.querySelector("footer"),
        };
      });

      fs.writeFileSync(
        path.join(outDir, `${target.name}-info.json`),
        JSON.stringify(info, null, 2)
      );

      await page.close().catch(() => {});
    } catch (err) {
      console.error(`Failed ${target.name}:`, err.message);
    }
  }

  await browser.close();
  console.log("\nRecon capture completed!");
}

run().catch(console.error);
