package utils

const cmName = "c7n-input-values"

/*
func GetInputValue(key string) string {
	//cm := getOrCreateCm(cmName)
	return cm.Data[key]
}

func SaveInputValue(key, value string) {
	if value != "" {
		//saveToCm(cmName, key, value)
	}
}

func saveToCm(cmName, key, value string) {
	kc := *context.Ctx.KubeClient

	cm := getOrCreateCm(cmName)
	cm.Data[key] = value
	_, err := kc.CoreV1().ConfigMaps(context.Ctx.Namespace).Update(cm)
	if err != nil {
		log.Error(err)
		os.Exit(122)
	}
}

func getOrCreateCm(cmName string) *v1.ConfigMap {
	context.Ctx.Mux.Lock()
	defer context.Ctx.Mux.Unlock()

	kc := *context.Ctx.KubeClient
	cm, err := kc.CoreV1().ConfigMaps(context.Ctx.Namespace).Get(cmName, meta_v1.GetOptions{})
	if errors.IsNotFound(err) {
		log.Info("creating configmap to cluster")
		data := make(map[string]string)
		data["user_info"] = fmt.Sprintf("email: %s", context.Ctx.Metrics.Mail)
		cm = &v1.ConfigMap{
			TypeMeta: meta_v1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: meta_v1.ObjectMeta{
				Name:   cmName,
				Labels: context.Ctx.CommonLabels,
			},
			Data: data,
		}
		cm, err = (*context.Ctx.KubeClient).CoreV1().ConfigMaps(context.Ctx.Namespace).Create(cm)
		if err != nil {
			context.Ctx.CheckExist(128, err.Error())
		}
	}
	return cm
}
*/
