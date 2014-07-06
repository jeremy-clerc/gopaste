#!/usr/bin/env perl

use strict;
use warnings;

my $curtime = time;
if ($#ARGV ne 0) {
    print "Usage: $0 path_to_data_dir/";
    exit 1;
}

print "Cleaning directory $ARGV[0] ...\n";
foreach my $file (glob $ARGV[0] . "/*") {
    if ($file =~ /$ARGV[0]\/([a-zA-Z0-9]+)-[a-z0-9]+-([0-9]+)/) {
        if ($2 < $curtime) {
            unlink $file;
            print "Cleaning paste $1\n";
        }
    }
}
