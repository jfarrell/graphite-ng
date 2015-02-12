/**
 * Structs are the basic complex data structures. They are comprised of fields
 * which each have an integer identifier, a type, a symbolic name, and an
 * optional default value.
 *
 * Fields can be declared "optional", which ensures they will not be included
 * in the serialized output if they aren't set.  Note that this requires some
 * manual management in some languages.
 */

# {"target": "test.metric1", "datapoints": [[0.000000, 1423683330]]}

#include "shared.thrift"

struct Datapoint {
  1: double value,
  2: i32 timestamp,
}

struct RenderData {
  1: string target
  2: list<Datapoint> datapoints
}

typedef list<string> MetricList

#service Calculator extends shared.SharedService {
service GraphiteNG {

  RenderData render()

  MetricList metrics()

}
