<script lang="ts">
    import { goto } from "$app/navigation";
    import { register } from "$lib/auth";
    import Button from "../../../components/form/Button.svelte";
    import Input from "../../../components/form/Input.svelte";
    import { toaster } from "../../../components/toaster/toaster";

    let email: string;
    let password: string;
    let invite: string;

    const handleRegister = async () => {
    await register(email, password, invite)
        .then(() => goto("/login"))
        .catch(() => toaster.push({ type: "error", message: "Failed to register." }, 5000));
    };

</script>

<div class="h-full w-full flex flex-col justify-center items-center">
    <form class="p-10 bg-white rounded-xl w-1/3" on:submit|preventDefault={handleRegister}>
        <Input bind:value={email} name="Email"/>
        <Input bind:value={password} name="Password" type="password"/>
        <Input bind:value={invite} name="Invite Code" placeholder="Registration currently restricted for users without an invite"/>
        <div class="flex flex-row justify-end">
            <Button type="submit">Register</Button>
        </div>
    </form>
</div>