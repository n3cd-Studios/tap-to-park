
<script>
  import { page } from '$app/stores';
  import { env } from '$env/dynamic/public';

  export let title = "Home | Example.com";
  export let description = "Description of your website.";

  let url = env.PUBLIC_FRONTEND_HOST ?? "https://0.0.0.0:8081";
  let image = "${url}/logo.png";
</script>

<svelte:head>
  <title>{title} | Tap-To-Park</title>
  <meta name="description" content={description}>
  <meta property="og_site_name" content="">
  <meta property="og:url" content="{url}{$page.url.pathname.toString()}">
  <meta property="og:type" content="website">
  <meta property="og:title" content="{title}">
  <meta property="og:description" content={description}>
  <meta property="og:image" content={image}>

  <meta name="twitter:card" content="summary_large_image">
  <meta property="twitter:domain" content="example.com">
  <meta property="twitter:url" content="{url}{$page.url.pathname.toString()}">
  <meta name="twitter:title" content="{title}">
  <meta name="twitter:description" content={description}>
  <meta name="twitter:image" content={image}>
  {@html `  <script type="application/ld+json">{
   "@context": "https://schema.org",
   "@type": "Website",
   "name": "${title} | example.com,
   "url": "${url}${$page.url.pathname}",
   "logo": ${image}  }</script>`}
</svelte:head>
