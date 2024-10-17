<script lang="ts">
    import moment from "moment";
    import AuthRoute from "../../../components/auth/AuthRoute.svelte";
    import Table from "../../../components/table/Table.svelte";
    import TableItem from "../../../components/table/TableItem.svelte";
    import Fa from "svelte-fa";
    import { faRemove } from "@fortawesome/free-solid-svg-icons";
    import { get } from "$lib/api";
    import { getAuthHeader } from "$lib/auth";
    import { toaster } from "../../../components/toaster/toaster";

    export let data;
    const { sessions } = data;

    const revokeSession = async (guid: string) => {
        await get({ route: `auth/sessions/${guid}`, headers: getAuthHeader(), method: "DELETE" });
        toaster.push({
            type: "success",
            message: "Revoked session."
        })
    }

</script>

<AuthRoute>
    <div class="m-10">
        <Table
            data={sessions}
            columns={["guid", "device", "ip", "expires", "last used", "revoke"]}
            loading={false}
            let:guid
            let:device
            let:ip
            let:expires
            let:lastUsed
        >
            <TableItem>{guid}</TableItem>
            <TableItem>{device}</TableItem>
            <TableItem>{ip}</TableItem>
            <TableItem>{moment(expires).fromNow()}</TableItem>
            <TableItem>{moment(lastUsed).fromNow()}</TableItem>
            <TableItem><button on:click={() => revokeSession(guid)}><Fa icon={faRemove} /></button></TableItem>
        </Table>
    </div>
</AuthRoute>
