<script lang="ts">
    import { goto } from "$app/navigation";
    import { login } from "$lib/auth";
    import Button from "../../../components/form/Button.svelte";
    import Input from "../../../components/form/Input.svelte";

    let email: string;
    let password: string;

    const handleLogin = async () => {
        await login(email, password);
        goto("/admin");
    };

</script>

<div class="h-full w-full flex flex-col justify-center items-center">
    <form class="p-10 bg-white rounded-xl w-1/3" on:submit|preventDefault={handleLogin}>
        <Input bind:value={email} name="Email"/>
        <Input bind:value={password} name="Password" type="password"/>
        <div class="flex flex-row justify-end">
            <Button type="submit">Login</Button>
        </div>
    </form>
</div>