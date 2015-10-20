require 'sshkit'
require 'sshkit/dsl'
require 'optparse'
require 'timeout'

OptionParser.new do |opts|
  opts.banner = "Usage: deploy.rb [options]"

  opts.on('-u', '--username NAME', 'Username') { |v| USER = v }
  opts.on('-s', '--server HOST', 'Destination host') { |v| SERVER = v }

end.parse!

APP_NAME = "imagenar-g"
LOCAL_TARGET = "bin"
APP_DIR = "/home/#{USER}/#{APP_NAME}"
DEST_TARGET = "#{APP_DIR}/releases/#{Time.now.strftime("%Y%m%dT%H%M%S")}"

run_locally do
  if test("[ -d #{LOCAL_TARGET} ]")
    execute :rm, '-rf', LOCAL_TARGET
  end
  execute :mkdir, "#{LOCAL_TARGET}"
  with 'GO15VENDOREXPERIMENT' => 1 do
    execute :go, :build, '-o', "#{LOCAL_TARGET}/#{APP_NAME}", "*.go"
  end
end

host = SSHKit::Host.new("#{USER}@#{SERVER}")
SSHKit::Backend::Netssh.config.pty = true

on host do
  execute :mkdir, '-p', "#{DEST_TARGET}/db"
  upload! "#{LOCAL_TARGET}/#{APP_NAME}", DEST_TARGET
  upload! 'db/migrations', "#{DEST_TARGET}/db", recursive: true
  execute :ln, '-s', "#{APP_DIR}/share/dbconf.yml", "#{DEST_TARGET}/db/dbconf.yml"
  execute :killall, APP_NAME rescue nil
  within APP_DIR do
    execute :ln, '-nfs', DEST_TARGET, :current
  end
  within "#{APP_DIR}/current" do
    # execute :goose, '-env', :production, :up
    # puts capture(:nohup, "#{APP_DIR}/current/#{APP_NAME}", "-e=production", "&")
    Timeout::timeout(1) {
      execute :nohup, "#{APP_DIR}/current/#{APP_NAME} -e=production < /dev/null > imagenar-g.log"
    } rescue nil
  end
end
