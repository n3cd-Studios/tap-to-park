
<script>
    import "../app.css";
    import "chart.js/auto";

    import Logo from "../assets/logo.png"
    import Toaster from "../components/toaster/Toaster.svelte";
    import UserDropdown from "../components/userDropdown/UserDropdown.svelte";
    import { page } from "$app/stores";
    
    let currentPage = "";
    $: {
        const path = $page.url.pathname;
        if (path.includes("/user")) {
            currentPage = "userPage";
        } else if (path.includes("/admin")) {
            currentPage = "adminPage";
        } else if (path.includes("/auth/login")) {
            currentPage = "loginPage";
        } else {
            currentPage = "defaultPage";
        }
    }
</script>

<div class="flex flex-col h-screen font-nunito">
    <div class="my-16 bg-white w-full flex justify-center items-center">
        <a href="/" aria-label="Go to homepage">
            <img src={Logo} alt="Logo" width={371} height={55}/>
        </a>
        <div class="absolute top-6 right-12">
            <UserDropdown {currentPage}  />
        </div>
    </div>
    <div class="h-full bg-primary rounded-t-[3rem]">
        <slot />
        <Toaster />
    </div>
</div>
