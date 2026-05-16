package leagueforbiddenequipment

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
	"starconflict/lib/variantdict"
)

func ScmdLeagueFrobiddenEquipment() ([]byte, error) {
	bw := bitwriter.NewWriter(make([]byte, 0, 500))
	if err := BwScmdLeagueFrobiddenEquipment(bw); err != nil {
		return nil, err
	}
	return bw.ReturnSlice(), nil
}

func BwScmdLeagueFrobiddenEquipment(bw *bitwriter.Writer) error {
	forbiddenEquipment := map[string]bool{
		"Weapon_Minigun_NY_T3":               true,
		"Weapon_Minigun_NY_T4":               true,
		"Weapon_Minigun_NY_T5":               true,
		"Weapon_SmartRocketLauncher_NY_T3":   true,
		"Weapon_SmartRocketLauncher_NY_T4":   true,
		"Weapon_SmartRocketLauncher_NY_T5":   true,
		"Weapon_SuppressingGun_NY_T3":        true,
		"Weapon_SuppressingGun_NY_T4":        true,
		"Weapon_SuppressingGun_NY_T5":        true,
		"FireStealingPot":                    true,
		"BushidoCall":                        true,
		"DronDamageMod":                      true,
		"EnemyInfoCheck":                     true,
		"CloakingDevice_Uniq_T3":             true,
		"FrigateDrone_T3_LessShldHeal":       true,
		"WolfRun":                            true,
		"SaveBack":                           true,
		"Module_Glass_PlasmaWeb":             true,
		"Module_Glass_PhaseShield":           true,
		"Module_Glass_EnergyShieldCommander": true,
		"Module_Glass_FrigateDrone":          true,
		"Module_Glass_CloakingDevice":        true,
		"Module_Glass_MicroWarp":             true,
		"Module_Glass_TimeBoost":             true,
	}
	return variantdict.BwMarshal(bw, forbiddenEquipment)
}

func SendScmdLeagueFrobiddenEquipment(conn net.Conn) {
	resp, err := ScmdLeagueFrobiddenEquipment()
	if err != nil {
		slog.Error("Failed to create ScmdLeagueFrobiddenEquipment()", "error", err)
	}
	conn.Write(protocol.MakeMessage(types.SCMD_LEAGUE_FORBIDDEN_EQUIPMENT, 0, 0, resp))
}
