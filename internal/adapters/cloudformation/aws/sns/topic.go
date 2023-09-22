package sns

import (
	"github.com/khulnasoft-lab/defsec/pkg/providers/aws/sns"
	"github.com/khulnasoft-lab/defsec/pkg/types"
	"github.com/khulnasoft-lab/vul-iac/pkg/scanners/cloudformation/parser"
)

func getTopics(ctx parser.FileContext) (topics []sns.Topic) {
	for _, r := range ctx.GetResourcesByType("AWS::SNS::Topic") {

		topic := sns.Topic{
			Metadata: r.Metadata(),
			ARN:      types.StringDefault("", r.Metadata()),
			Encryption: sns.Encryption{
				Metadata: r.Metadata(),
				KMSKeyID: r.GetStringProperty("KmsMasterKeyId"),
			},
		}

		topics = append(topics, topic)
	}
	return topics
}
