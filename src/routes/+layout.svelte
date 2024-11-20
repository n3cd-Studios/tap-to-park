
<script>
    import "../app.css";
    import "chart.js/auto";

    import Logo from "../assets/logo.png"
    import Toaster from "../components/toaster/Toaster.svelte";
</script>

<div class="flex flex-col h-screen font-nunito">
    <div class="my-16 bg-white w-full flex justify-center items-center">
        <a href="/" aria-label="Go to homepage">
            <img src={Logo} alt="Logo" width={371} height={55}/>
        </a>
    </div>
    <div class="h-full bg-primary rounded-t-[3rem]">
        <slot />
        <Toaster />
    </div>
</div>
