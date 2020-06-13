


func main() { 
    router := httprouter.New() 
    router.GET("/", controller.Index) 
    router.POST("/filter", controller.Filter) 
    router.POST("/prioritize", controller.Prioritize) 
    log.Printf("start up sample-scheduler-extender!\n") 
    log.Fatal(http.ListenAndServe(":8888", router)) 
}

func filter(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult { 
    var filteredNodes []v1.Node 
    failedNodes := make(schedulerapi.FailedNodesMap) 
    pod := args.Pod 
    for _, node := range args.Nodes.Items { 
        fits, failReasons, _ := podFitsOnNode(pod, node) 
        if fits { 
            filteredNodes = append(filteredNodes, node) 
        } else { 
            failedNodes[node.Name] = strings.Join(failReasons, ",") 
        } 
    }
    result := schedulerapi.ExtenderFilterResult{ 
        Nodes: &v1.NodeList{ 
            Items: filteredNodes, 
        },
        FailedNodes: failedNodes, 
        Error: "", 
    }
    return &result
}

func podFitsOnNode(pod *v1.Pod, node v1.Node) (bool, []string, error) { 
    fits := true 
    var failReasons []string 
    for _, predicateKey := range predicatesSorted { 
        fit, failures, err := predicatesFuncs[predicateKey](pod, node) 
        if err != nil { 
            return false, nil, err 
        }
        fits = fits && fit 
        failReasons = append(failReasons, failures...) 
    }
    return fits, failReasons, nil 
}
func LuckyPredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) { 
    lucky := strings.Contains(pod.Name,"87")
    if lucky { 
        log.Printf("Pod %v/%v is lucky to fit on node %v\n", pod.Name, pod.Namespace, node.Name) 
        return true, nil, nil 
    }
    log.Printf("pod %v/%v is unlucky to fit on node %v\n", pod.Name, pod.Namespace, node.Name) 
    return false, []string{LuckyPredFailMsg}, nil 
}
func prioritize(args schedulerapi.ExtenderArgs) *schedulerapi.HostPriorityList { 
    pod := args.Pod 
    nodes := args.Nodes.Items 
    hostPriorityList := make(schedulerapi.HostPriorityList, len(nodes)) 
    for i, node := range nodes { 
        score := len(pod.Name) % (schedulerapi.MaxPriority + 1) 
        log.Printf(luckyPrioMsg, pod.Name, pod.Namespace, score) 
        hostPriorityList[i] = schedulerapi.HostPriority{ 
            Host: node.Name, Score: score, 
        } 
    }
    return &hostPriorityList 
}
