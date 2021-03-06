package targets

import (
	"cosmos-validator-mission-control/config"
	"encoding/json"
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
)

// CheckGaiad function to get gaiad status and send
//alerts to telgram and email accounts
func CheckGaiad(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Println("Error while getting gaiad status..")
		return
	}

	var status ValidatorRpcStatus
	err = json.Unmarshal(resp.Body, &status)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	caughtUp := !status.Result.SyncInfo.CatchingUp
	if !caughtUp {
		_ = SendTelegramAlert(fmt.Sprintf("Gaiad on your validator instance is not running: \n%v", string(resp.Body)), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Gaiad on your validator instance is not running: \n%v", string(resp.Body)), cfg)
		_ = writeToInfluxDb(c, bp, "vcf_gaiad_status", map[string]string{}, map[string]interface{}{"status": 0})
		return
	}

	_ = writeToInfluxDb(c, bp, "vcf_gaiad_status", map[string]string{}, map[string]interface{}{"status": 1})
}
