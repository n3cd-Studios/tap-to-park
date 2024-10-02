<script lang="ts">
    import { goto } from "$app/navigation";
    import { login } from "$lib/auth";
    import Button from "../../../components/Button.svelte";
    import Input from "../../../components/Input.svelte";

    let email: string;
    let password: string;

</script>

<div class="h-full w-full flex flex-col justify-center items-center">
    <div class="p-10 bg-white rounded-xl w-1/3">
        <Input bind:value={email} name="Email" />
        <Input bind:value={password} name="Password" type="password"/>
        <div class="flex flex-row justify-end">
            <Button click={async () => {
                await login(email, password);
                goto("/admin");
            }}>Login</Button>
        </div>
    </div>
</div>