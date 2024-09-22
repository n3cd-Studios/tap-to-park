
<script>
    import "../app.css";
    import Logo from "../assets/logo.png"
</script>

<div class="flex h-screen flex-col">
    <div class="h-1/3 w-full flex justify-center items-center">
        <img src={Logo} />
    </div>
    <div class="h-full bg-primary rounded-t-[3rem]">
        <slot />
    </div>
</div>