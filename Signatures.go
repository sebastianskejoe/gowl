package gowl

type Signature struct {
    opcode int16
    signature string
}

var signatures map[string]Signature

func init() {
    signatures = make(map[string]Signature, 0)
	signatures["wl_display_sync"] = Signature{0, "n"}
	signatures["wl_display_get_registry"] = Signature{1, "n"}

	signatures["wl_registry_bind"] = Signature{0, "usun"}


	signatures["wl_compositor_create_surface"] = Signature{0, "n"}
	signatures["wl_compositor_create_region"] = Signature{1, "n"}

	signatures["wl_shm_pool_create_buffer"] = Signature{0, "niiiiu"}
	signatures["wl_shm_pool_destroy"] = Signature{1, ""}
	signatures["wl_shm_pool_resize"] = Signature{2, "i"}

	signatures["wl_shm_create_pool"] = Signature{0, "nhi"}

	signatures["wl_buffer_destroy"] = Signature{0, ""}

	signatures["wl_data_offer_accept"] = Signature{0, "us"}
	signatures["wl_data_offer_receive"] = Signature{1, "sh"}
	signatures["wl_data_offer_destroy"] = Signature{2, ""}

	signatures["wl_data_source_offer"] = Signature{0, "s"}
	signatures["wl_data_source_destroy"] = Signature{1, ""}

	signatures["wl_data_device_start_drag"] = Signature{0, "ooou"}
	signatures["wl_data_device_set_selection"] = Signature{1, "ou"}

	signatures["wl_data_device_manager_create_data_source"] = Signature{0, "n"}
	signatures["wl_data_device_manager_get_data_device"] = Signature{1, "no"}

	signatures["wl_shell_get_shell_surface"] = Signature{0, "no"}

	signatures["wl_shell_surface_pong"] = Signature{0, "u"}
	signatures["wl_shell_surface_move"] = Signature{1, "ou"}
	signatures["wl_shell_surface_resize"] = Signature{2, "ouu"}
	signatures["wl_shell_surface_set_toplevel"] = Signature{3, ""}
	signatures["wl_shell_surface_set_transient"] = Signature{4, "oiiu"}
	signatures["wl_shell_surface_set_fullscreen"] = Signature{5, "uuo"}
	signatures["wl_shell_surface_set_popup"] = Signature{6, "ouoiiu"}
	signatures["wl_shell_surface_set_maximized"] = Signature{7, "o"}
	signatures["wl_shell_surface_set_title"] = Signature{8, "s"}
	signatures["wl_shell_surface_set_class"] = Signature{9, "s"}

	signatures["wl_surface_destroy"] = Signature{0, ""}
	signatures["wl_surface_attach"] = Signature{1, "oii"}
	signatures["wl_surface_damage"] = Signature{2, "iiii"}
	signatures["wl_surface_frame"] = Signature{3, "n"}
	signatures["wl_surface_set_opaque_region"] = Signature{4, "o"}
	signatures["wl_surface_set_input_region"] = Signature{5, "o"}
	signatures["wl_surface_commit"] = Signature{6, ""}

	signatures["wl_seat_get_pointer"] = Signature{0, "n"}
	signatures["wl_seat_get_keyboard"] = Signature{1, "n"}
	signatures["wl_seat_get_touch"] = Signature{2, "n"}

	signatures["wl_pointer_set_cursor"] = Signature{0, "uoii"}




	signatures["wl_region_destroy"] = Signature{0, ""}
	signatures["wl_region_add"] = Signature{1, "iiii"}
	signatures["wl_region_subtract"] = Signature{2, "iiii"}

}