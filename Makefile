INSTALL = install
prefix  ?= /usr/local
SHELL := /bin/bash

PROGRAMS := precheck copy service
	
$(phony all): $(PROGRAMS)

.PHONY install: precheck 
	

DESTDIR?=/usr/local/simpcord

precheck:
	@echo "Checking if simpcord is already installed"
	@if [ -f $(DESTDIR)/simpcord ]; then \
		echo -e "\033[34mAlready Installed!\e[0m Performing upgrade"; \
		$(MAKE) upgrade; \
	else \
		$(MAKE) copy; \
		$(MAKE) service; \
	fi
	
copy:
	@echo "Copying files";
	@test -d $(prefix) || mkdir -p $(prefix); 
	@$(INSTALL) -d -m 755 $(prefix)/simpcord;
	@$(INSTALL) -m 755 simpcord $(prefix)/simpcord/;

upgrade:
	@systemctl stop simpcord.service; 
	@$(INSTALL) -m 755 simpcord $(prefix)/simpcord/;
	@systemctl start simpcord.service; 

uninstall:
	@if [ ! -f $(DESTDIR)/simpcord ]; then \
		echo "simpcord is not installed!"; \
		exit 1; \
	fi
	@echo "I am about to remove $(DESTDIR) and all of its subdirectories"
	@read -r -p "Are you sure? [y/N] " response; case "$$response" in [yY][eE][sS]|[yY]) \
		if [ -f /etc/redhat-release ]; then \
			systemctl stop simpcord; \
			systemctl disable simpcord; \
			rm /etc/systemd/system/simpcord.service; \
			rm /etc/sysconfig/simpcord; \
			systemctl daemon-reload; \
			systemctl reset-failed; \
		elif [ -f /etc/debian_version ]; then \
			if [ $$(ps --no-headers -o comm 1) = "systemd" ]; then \
				systemctl stop simpcord; \
				systemctl disable simpcord; \
			fi ; \
			rm /etc/systemd/system/simpcord.service; \
			rm /etc/default/simpcord; \
			if [ $$(ps --no-headers -o comm 1) = "systemd" ]; then \
				systemctl daemon-reload; \
				systemctl reset-failed; \
			fi ; \
		fi ; \
		rm -rf $(DESTDIR);  \
	;; \
	esac

service:
	@echo "Installing simpcord Service..."
	@if [ $$(ps --no-headers -o comm 1) = "systemd" ]; then \
		if [ -f /etc/redhat-release ]; then \
			$(INSTALL) -m 755 simpcord.service /etc/systemd/system/simpcord.service; \
			$(INSTALL) -m 640 simpcord.options /etc/sysconfig/simpcord; \
		elif [ -f /etc/debian_version ]; then \
			$(INSTALL) -m 755 simpcord.service.debian /etc/systemd/system/simpcord.service; \
			$(INSTALL) -m 640 simpcord.options /etc/default/simpcord; \
		fi ; \
		systemctl daemon-reload; \
		systemctl enable simpcord.service; \
		systemctl start simpcord.service; \
	else \
	    echo "\033[1;33msimpcord service requires systemd\033[0m: Could not install service unit file"; \
	fi; \
