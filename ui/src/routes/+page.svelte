<script lang="ts">
	import * as VIAM from '@viamrobotics/sdk';
	import { Button, SearchableSelect } from '@viamrobotics/prime-core';

	{/*TODO: 	
		* Remove function into its own file
		* Modify buf.yaml to generate typescript code
		* Implement routes fetch from proto
		* Call start navigation
	*/}
	let loadRobot = connectToRobot();

	async function connectToRobot() {
		const host = 'fleet-cam-01-main.ve4ba7w5qr.viam.cloud';

		const robot = await VIAM.createRobotClient({
			host,
			credential: {
				type: 'robot-location-secret',
				payload: 'emqt8rz43ffblxw4l4mj5j2cobmqiuxccw4w1hsx0ny8m86t'
			},
			authEntity: host,
			signalingAddress: 'https://app.viam.com:443'
		});

		// camera
		const cameraClient = new VIAM.CameraClient(robot, 'camera');
		const cameraImage = await cameraClient.getImage();

		if (!cameraImage) return 'Not OK';

		return 'OK';
	}
</script>

<svelte:head>
	<title>Robot control</title>
	<meta name="description" content="Robot remote control" />
</svelte:head>

{#await loadRobot}
	Connecting to robot...
{:then value}
	{#if value === 'OK'}
		<div class="flex gap-8 m-auto">
			<!-- TODO: dynamically fetch routes -->
			<div class="text-bold">Routes:</div>
			<SearchableSelect options={['Route 1', 'Route 2']} />
		</div>

		<div class="flex gap-12 justify-center">
			<Button variant="success" on:click={() => (loadRobot = connectToRobot())}>
				Start robot navigation
			</Button>
			<Button variant="danger" on:click={() => {}}>Stop robot navigation</Button>
		</div>
	{:else}
		<div class="flex justify-center mt-4">
			<Button icon="refresh" variant="dark" on:click={() => (loadRobot = connectToRobot())}>
				Retry robot connection
			</Button>
		</div>
	{/if}
{:catch someError}
	System error: {someError.message}.
{/await}
