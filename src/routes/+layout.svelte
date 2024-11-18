
<script>
    import "../app.css";
    import Logo from "../assets/logo.png"
    import Toaster from "../components/toaster/Toaster.svelte";
</script>

<div class="flex h-screen flex-col font-nunito">
    <div class="h-1/3 w-full flex justify-center items-center">
        <a href="/" aria-label="Go to homepage">
            <img src={Logo} alt="Logo" width={371} height={55}/>
        </a>
    </div>
    <div class="h-full bg-primary rounded-t-[3rem]">
        <slot />
        <Toaster />
    </div>
</div>