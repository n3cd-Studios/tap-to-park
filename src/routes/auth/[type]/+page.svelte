
<script lang="ts">
    import { goto } from "$app/navigation";
    import { authStore, getUserInfo } from "$lib/auth";
    import { UserRole } from "$lib/models";
    import { onMount } from "svelte";

    export let data;

    onMount(async () => {
        const res = await fetch(`http://localhost:8080/api/auth/${data.type}/callback?${data.search.toString()}`, { method: "GET" })
            .then(res => res.json())
            .catch(_ => null);

        if (res) $authStore.token =  res.token;
        
        const user = await getUserInfo();
        if (!user) return;
        if (user.role == UserRole.ADMIN) goto("/admin")
        else goto("/user");
    })
</script>

<div class="h-full flex flex-row justify-center items-center">
    <p class="text-xl">Give us a second while we sign you in..</p>
</div>