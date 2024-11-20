<script lang="ts">
    import { Paginator } from "$lib/api";
    import { getAuthHeader } from "$lib/auth";
    import { type Spot, type Reservation } from "$lib/models";
    import { onMount } from "svelte";
    import Table from "../../../components/table/Table.svelte";
    import TableItem from "../../../components/table/TableItem.svelte";
    import Button from "../../../components/form/Button.svelte";
    import { Formats, pluralize } from "$lib/lang";
    import moment from "moment";
    import Chart from "../../../components/Chart.svelte";
    import { getWithDefault } from "$lib/api";

    let loading = false;

    let items: Reservation[] = [];
    let paginator = new Paginator<Reservation>(
        {
            route: "organization/reservations",
            headers: getAuthHeader(),
            method: "GET",
        },
        10,
    );
    paginator.subscribe((cb) => (items = cb));

    type TopSpot = { name: string; id: string; revenue: number };
    type PeakTime = { amount: number, time: number, revenue: number };
    type RevenueByMonth = { month: number, revenue: number };

    let topSpots: TopSpot[] = [];
    let peakTimes: PeakTime[] = [];
    let revenueByMonth: RevenueByMonth[] = [];

    onMount(async () => {
        await paginator.load();
        topSpots = await getWithDefault<TopSpot[]>({ route: "analytics/top", headers: getAuthHeader() }, []);
        peakTimes = await getWithDefault<PeakTime[]>({ route: "analytics/peak", headers: getAuthHeader() }, []);
        revenueByMonth = await getWithDefault<RevenueByMonth[]>({ route: "analytics/revenue", headers: getAuthHeader() }, []);
    });
</script>

<div class="flex flex-row gap-2">
    <div class="bg-white w-full rounded-lg p-3">
        <Chart
            config={{
                type: "bar",
                data: {
                    labels: revenueByMonth.map((revenue) => Formats.Month(revenue.month - 1)),
                    datasets: [
                        {
                            label: "Revenue",
                            data: revenueByMonth.map((revenue) => revenue.revenue),
                            borderWidth: 1,
                        },
                    ],
                }
            }}
        />
    </div>
    <div class="bg-white w-full rounded-lg p-3">
        <Chart
            config={{
                type: "bar",
                data: {
                    labels: topSpots.map((spot) => spot.name),
                    datasets: [
                        {
                            label: "Revenue",
                            data: topSpots.map((spot) => spot.revenue),
                            borderWidth: 1,
                        },
                    ],
                }
            }}
        />
    </div>
    <div class="bg-white w-full rounded-lg p-3">
        <Chart
            config={{
                type: "bar",
                data: {
                    labels: peakTimes.map((time) => Formats.Time(time.time)),
                    datasets: [
                        {
                            label: "Transactions",
                            data: peakTimes.map((time) => time.amount),
                            borderWidth: 1,
                        },
                        {
                            label: "Revenue",
                            data: peakTimes.map((time) => time.revenue / 100),
                            borderWidth: 1,
                        },
                    ],
                },
            }}
        />
    </div>
</div>
<div class="flex flex-col gap-2 mt-2">
    <Table
        columns={["email", "start", "end", "duration"]}
        data={items}
        {loading}
        let:email
        let:start
        let:end
        let:price
    >
        <TableItem>{email}</TableItem>
        <TableItem>{Formats.Date(start)}</TableItem>
        <TableItem>{Formats.Date(end)}</TableItem>
        <TableItem
            >{pluralize(
                moment(moment(start).diff(moment(end))).hours(),
                "hour",
            )} for
            {Formats.USDollar.format(price / 100)}</TableItem
        >
    </Table>
    <div class="flex flex-row justify-center gap-2">
        <Button
            on:click={() => paginator.last()}
            aria-label="Go to previous page">{"<"}</Button
        >
        <Button on:click={() => paginator.next()} aria-label="Go to next page"
            >{">"}</Button
        >
    </div>
</div>
