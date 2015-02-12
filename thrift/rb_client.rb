$:.push('gen-rb')
#$:.unshift '../../lib/rb/lib'

require 'thrift'

require 'graphite_n_g'

begin
  endpoint = ARGV[0] || "metrics"

  transport = Thrift::BufferedTransport.new(Thrift::Socket.new('127.0.0.1', 9090))
  protocol = Thrift::JsonProtocol.new(transport)
  client = GraphiteNG::Client.new(protocol)

  transport.open()

  case endpoint
  when "metrics"
    data = client.metrics()
    puts "Metrics:"
    data.each do |d|
      puts "#{d}"
    end
  when "render"
    data = client.render()
    puts "Target: #{data.target}\nDatapoints:"
    data.datapoints.each do |d|
      puts "Timestamp: #{d.timestamp} - Value #{d.value}"
    end
  end

  transport.close()

rescue Thrift::Exception => tx
  print 'Thrift::Exception: ', tx.message, "\n"
end
