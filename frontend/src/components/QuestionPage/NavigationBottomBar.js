// TODO: style navigation bottom bar
import * as React from 'react';
import { Box, Button } from "@mui/material";

function NavigationBottomBar() {
	return (
		<>
			<Box display="flex" width="100%" justifyContent="flex-end"><Button>Save</Button></Box>
			<Box display="flex" justifyContent="space-between">
				<Box display="flex" width="100%" justifyContent="center">
					<Button>Prev</Button>
					<Button>Next</Button>
				</Box>
			</Box>
		</>
	);
}

export default NavigationBottomBar