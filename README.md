# scheduler-extender

学号：17341187

姓名：杨勰



kubernetes在master节点上用kubelet检测容器的状态，即Pod信息中Spec.nodeName字段的变化，然后将Pod进入预选阶段，判断pod.Name中是否带“87”字段，该阶段通过一系列的预选算法选出集群中适合Pod运行的节点，带着这些信息进入Priorities阶段。该阶段通过Priorities算法为适合该Pod调度对每个Node进行打分，最后选出集群中分数最高的Pod运行的一个节点，从而实现Pod的调度。
